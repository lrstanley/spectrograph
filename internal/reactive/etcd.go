// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package reactive

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/semver"
	"github.com/apex/log"
	"github.com/denisbrodbeck/machineid"
	"github.com/lrstanley/spectrograph/internal/models"
	etcdcli "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// https://github.com/etcd-io/etcd/blob/master/tests/integration/clientv3/concurrency/example_election_test.go
// https://medium.com/@felipedutratine/distributed-lock-with-etcd-in-go-d21e7df145bc
// https://github.com/kubernetes/client-go/blob/master/examples/leader-election/main.go
// https://github.com/Masterminds/semver

const (
	reconnectThrottle = 5 * time.Second
	dialTimeout       = 10 * time.Second
	leaseSeconds      = 15
)

var (
	ErrIncompatibleVersion = errors.New("current etcd leader is running an incompatible version")
)

type Campaign struct {
	cli         models.EtcdConfig
	logger      log.Interface
	wantsLeader bool
	node        string
	version     string

	etcd          *etcdcli.Client
	LeaderUpdates chan *LeaderInfo

	// Mutex is only used to protect single-threaded writes from breaking
	// multi-threaded reads from external clients of this package. Since we
	// know there are no other writes other than for setup, we don't have to
	// protect internal reads.
	mu       sync.RWMutex
	session  *concurrency.Session
	election *concurrency.Election
}

// NewElection initiates a new campaign, that connects to etcd and campaigns
// for leader.
func NewElection(logger log.Interface, cli models.EtcdConfig, wantsLeader bool, version string) (campaign *Campaign) {
	entry := logger.WithFields(log.Fields{
		"source":         "etcd",
		"username":       cli.Username,
		"etcd-endpoints": strings.Join(cli.Endpoints, ","),
	})

	c := &Campaign{
		cli:           cli,
		logger:        entry,
		wantsLeader:   wantsLeader,
		version:       version,
		LeaderUpdates: make(chan *LeaderInfo),
	}

	var err error

	// Try to get the machine-id first, if it is available, fallback to hostname.
	c.node, err = machineid.ID()
	if err != nil {
		c.logger.WithError(err).Debug("unable to obtain machine-id")

		c.node, err = os.Hostname()
		if err != nil {
			c.logger.WithError(err).Debug("unable to obtain hostname (as machine-id)")
		}
	}

	if c.node == "" {
		c.logger.Fatal("unable to automatically identify node name")
	}

	c.logger = c.logger.WithField("node", c.node)

	c.logger.Info("machine-id acquired")
	c.logger.Info("initializing etcd connections")
	c.etcd, err = etcdcli.New(etcdcli.Config{
		Endpoints:   cli.Endpoints,
		DialTimeout: dialTimeout,

		Username: cli.Username,
		Password: cli.Password,
	})
	if err != nil {
		c.logger.WithError(err).Fatal("error initializing connection to etcd")
	}

	return c
}

func (c *Campaign) Close() error {
	return c.etcd.Close()
}

func (c *Campaign) IsLeader(ctx context.Context) (leaderInfo *LeaderInfo, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.election == nil {
		return nil, nil
	}
	resp, err := c.election.Leader(ctx)
	if err != nil {
		if err != concurrency.ErrElectionNoLeader {
			return nil, err
		}
		return nil, nil
	}

	value := string(resp.Kvs[0].Value)
	return newLeaderInfo(c, value), nil
}

func (c *Campaign) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			c.LeaderElection(ctx)
			select {
			case <-time.After(reconnectThrottle):
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (c *Campaign) LeaderElection(ctx context.Context) {
	var err error

	c.mu.Lock()
	c.session, err = concurrency.NewSession(c.etcd, concurrency.WithTTL(leaseSeconds))
	if err != nil {
		c.logger.WithError(err).Error("unable to establish new etcd election session")
		c.mu.Unlock()
		return
	}
	defer c.session.Close()

	c.election = concurrency.NewElection(c.session, "/leader/")
	observations := c.election.Observe(ctx)
	c.mu.Unlock()

	select {
	case <-ctx.Done():
		c.logger.Debug("context closed")
		return
	default:
	}

	lctx, lcancel := context.WithCancel(ctx)
	defer lcancel()

	wg := &sync.WaitGroup{}

	if c.wantsLeader {
		// Only campaign if we want to be leader. This is useful for workers.
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.logger.Debug("requesting election")
			if err = c.election.Campaign(lctx, fmt.Sprintf("%s=%s", c.node, c.version)); err != nil {
				c.logger.WithError(err).Error("error initiating campaign to become leader")
				lcancel()
			}
		}()
	}

	for {
		select {
		case <-lctx.Done():
			c.logger.Debug("campaign context closed, cleaning up")
			goto cleanup
		case <-ctx.Done():
			// We want to resign because we're shutting down.
			c.logger.Debug("context closed, cleaning up")
			goto cleanup
		case <-c.session.Done():
			// We've lost leader/expired/etc.
			c.logger.Debug("session closed, cleaning up")
			goto cleanup
		case resp, ok := <-observations:
			if !ok {
				c.logger.Debug("observations channel closed, cleaning up")
				goto cleanup
			}

			value := string(resp.Kvs[0].Value)
			leaderInfo := newLeaderInfo(c, value)

			if leaderInfo.Leader() {
				c.logger.Debug("acquired leader")
			}

			c.LeaderUpdates <- leaderInfo
		}
	}

cleanup:
	lcancel() // Tell campaign to stop running.
	cctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = c.election.Resign(cctx); err != nil {
		c.logger.WithError(err).Error("unable to resign from leader")
	}

	// Wait for any goroutines to exit.
	wg.Wait()
}

func (c *Campaign) WaitForLeader(ctx context.Context, checkTimeout, checkDelay time.Duration) (leader *LeaderInfo) {
	lctx, lcancel := context.WithTimeout(ctx, checkTimeout)
	defer lcancel()

	var err error
	for {
		leader, err = c.IsLeader(lctx)
		if err != nil {
			c.logger.WithError(err).Error("unable to obtain leader information")
		} else if leader != nil {
			// We're the leader, we can proceed.
			if leader.Leader() {
				return
			}

			if leader.IsCompatibleVersion() {
				// We're not leader, but the leader is a compatible version.
				return
			} else {
				c.logger.WithError(ErrIncompatibleVersion).Error("waiting for new leader")
			}
		}

		select {
		case <-time.After(checkDelay):
		case <-lctx.Done():
			c.logger.Error("context timed out while waiting for new leader")
			return
		}
	}
}

func versionSplit(input string) (node, version string) {
	split := strings.SplitN(input, "=", 2)
	if len(split) != 2 {
		panic(fmt.Sprintf("supplied version %q does not match node=version format", input))
	}

	return split[0], split[1]
}

type LeaderInfo struct {
	Node    string // Leader node name.
	Version string // Leader node version.
	Value   string // Leader node raw value.

	campaign *Campaign
}

func newLeaderInfo(c *Campaign, value string) *LeaderInfo {
	node, version := versionSplit(value)

	return &LeaderInfo{
		Node:     node,
		Version:  version,
		Value:    value,
		campaign: c,
	}
}

// Leader returns true if we're the active leader.
func (l *LeaderInfo) Leader() bool {
	return l.campaign.node == l.Node
}

func (l *LeaderInfo) IsCompatibleVersion() bool {
	if l.Version == "master" || l.campaign.version == "master" {
		if l.Version == "master" && l.campaign.version == "master" {
			return true
		}
		return false
	}

	var err error
	leaderVersion, err := semver.NewVersion(l.Version)
	if err != nil {
		l.campaign.logger.WithError(err).WithFields(log.Fields{
			"node":    l.Node,
			"version": l.Version,
		}).Debug("unable to parse leader version")
		return false
	}

	srcVersion, err := semver.NewVersion(l.campaign.version)
	if err != nil {
		l.campaign.logger.WithError(err).WithFields(log.Fields{
			"node":    l.campaign.node,
			"version": l.campaign.version,
		}).Debug("unable to parse source version")
		return false
	}

	constraint, err := semver.NewConstraint(fmt.Sprintf("~%d.%d.x-0", srcVersion.Major(), srcVersion.Minor()))
	if err != nil {
		l.campaign.logger.WithError(err).WithField("version", l.campaign.version).Fatal("unable to generate a constraint from current version")
	}

	ok, _ := constraint.Validate(leaderVersion)
	return ok
}
