// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"github.com/go-chi/chi"
)

func registerFetchRoutes(r chi.Router) {
	// r.Get("/api/v1/fetch-test/{providerkey:[a-zA-Z0-9._-]+}/*", func(w http.ResponseWriter, r *http.Request) {
	// 	provider, err := svcProviders.GetByKey(r.Context(), chi.URLParam(r, "providerkey"))
	// 	if helpers.HTTPError(w, r, http.StatusInternalServerError, err) {
	// 		return
	// 	}

	// 	proj := &models.Project{
	// 		ID:                primitive.NewObjectID(),
	// 		CreatedAt:         time.Now(),
	// 		UpdatedAt:         time.Now(),
	// 		Provider:          provider.ID,
	// 		ProviderArguments: make(map[string]string),
	// 	}

	// 	err = proj.ParseRawArguments(provider, strings.Split(chi.URLParam(r, "*"), "/"))
	// 	if helpers.HTTPError(w, r, http.StatusBadRequest, err) {
	// 		return
	// 	}

	// 	uri, err := provider.GenerateURI(proj.ProviderArguments)
	// 	if helpers.HTTPError(w, r, http.StatusInternalServerError, err) {
	// 		return
	// 	}

	// 	var foundProject *models.Project
	// 	foundProject, err = svcProject.GetByProvider(r.Context(), proj.Provider.Hex(), proj.ProviderArguments)
	// 	if err != nil && models.IsNotFound(err) {
	// 		err = svcProject.Create(r.Context(), proj)
	// 	}
	// 	if helpers.HTTPError(w, r, http.StatusInternalServerError, err) {
	// 		return
	// 	}
	// 	if foundProject != nil {
	// 		proj = foundProject
	// 	}

	// 	if proj.Disabled {
	// 		helpers.HTTPError(w, r, http.StatusNotFound, errors.New("this project has been disabled"))
	// 		return
	// 	}

	// 	if proj.Deleted {
	// 		helpers.HTTPError(w, r, http.StatusNotFound, errors.New("this project is pending deletion"))
	// 		return
	// 	}

	// 	// TODO: include request id, and enforce no dups based on request id?
	// 	req := &models.ScanRequest{
	// 		ID:        primitive.NewObjectID(),
	// 		ProjectID: proj.ID.Hex(),
	// 		URI:       uri,
	// 		Reference: "refs/heads/master",
	// 		Submitted: time.Now(),
	// 	}

	// 	var latestScan *models.ProjectScanResult
	// 	latestScan, err = svcProject.GetLatestScan(r.Context(), proj.ID.Hex())
	// 	if err != nil && !models.IsNotFound(err) {
	// 		helpers.HTTPError(w, r, http.StatusInternalServerError, err)
	// 		return
	// 	}

	// 	if latestScan != nil { //  && time.Since(latestScan.CompletedTimestamp) < 1*time.Hour
	// 		// TODO: if greater than an hour, still return it, but kick off a new scan?
	// 		w.WriteHeader(http.StatusOK)
	// 		pt.JSON(w, r, pt.M{"data": pt.M{
	// 			"project":     proj,
	// 			"latest_scan": latestScan,
	// 		}})
	// 		return
	// 	}

	// 	// TODO: IF the request leads to it being in the queue, optionally
	// 	// support waiting for a certain amount of time (with timeout) to hear
	// 	// a response back that the scan finished (and the results).

	// 	// Add it to the DB before we pass it in to be scanned. This may have
	// 	// a useful side effect in that it will let us know if there is already
	// 	// a queued scan for this project.
	// 	req.Updated = time.Now()
	// 	err = svcScanRequests.Create(r.Context(), req)
	// 	if models.IsDuplicate(err) {
	// 		// Assume that it was already in the queue.
	// 		w.WriteHeader(http.StatusAccepted)
	// 		pt.JSON(w, r, pt.M{"queued": true})
	// 		return
	// 	}
	// 	if helpers.HTTPError(w, r, http.StatusInternalServerError, err) {
	// 		return
	// 	}

	// 	select {
	// 	case scanRequestChan <- req:
	// 		w.WriteHeader(http.StatusAccepted)
	// 		pt.JSON(w, r, pt.M{"queued": true})
	// 	default:
	// 		helpers.HTTPError(w, r, http.StatusServiceUnavailable, errors.New("failed to queue project scan: server busy or unavailable"))
	// 		return
	// 	}
	// })
}
