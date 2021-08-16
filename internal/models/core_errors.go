// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

type ErrNotFound struct{ Err error }

func (e ErrNotFound) Error() string { return e.Err.Error() }

type ErrDuplicate struct{ Err error }

func (e ErrDuplicate) Error() string { return e.Err.Error() }

type ErrClientError struct{ Err error }

func (e ErrClientError) Error() string { return e.Err.Error() }

type ErrValidationFailed struct{ Err error }

func (e ErrValidationFailed) Error() string { return e.Err.Error() }

func IsClientError(e error) bool {
	if e == nil {
		return false
	}

	if _, ok := e.(*ErrClientError); ok {
		return true
	}
	return false
}

func IsNotFound(e error) bool {
	if e == nil {
		return false
	}

	if cerr, ok := e.(*ErrClientError); ok {
		if _, ok := cerr.Err.(*ErrNotFound); ok {
			return true
		}
	}
	return false
}

func IsDuplicate(e error) bool {
	if e == nil {
		return false
	}

	if cerr, ok := e.(*ErrClientError); ok {
		if _, ok := cerr.Err.(*ErrDuplicate); ok {
			return true
		}
	}
	return false
}
