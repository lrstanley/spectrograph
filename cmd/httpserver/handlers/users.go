// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/chix/xauth/v2"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/markbates/goth/gothic"
)

func (h *handler) getSelf(w http.ResponseWriter, r *http.Request) {
	chix.JSON(w, r, http.StatusOK, xauth.IdentFromContext[ent.User](r.Context()))
}

func (h *handler) deleteSelf(w http.ResponseWriter, r *http.Request) {
	self := xauth.IdentFromContext[ent.User](r.Context())

	if err := h.cfg.DB.User.DeleteOneID(self.ID).Exec(r.Context()); err != nil {
		chix.Error(w, r, err)
		return
	}

	_ = gothic.Logout(w, r)
	w.WriteHeader(http.StatusNoContent)
}

type banUserBody struct {
	Reason string `json:"reason" validate:"required,min=1,max=500"`
}

// userIDParam returns the target user ID from the URL, responding with an error
// if it's invalid.
func userIDParam(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		chix.ErrorWithCode(w, r, http.StatusBadRequest, errors.New("invalid user ID"))
		return 0, false
	}
	return id, true
}

func (h *handler) postBanUser(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDParam(w, r)
	if !ok {
		return
	}

	var body banUserBody
	if err := chix.Bind(r, &body); err != nil {
		chix.Error(w, r, err)
		return
	}

	self := xauth.IdentFromContext[ent.User](r.Context())
	if id == self.ID {
		chix.ErrorWithCode(w, r, http.StatusBadRequest, errors.New("cannot ban yourself"))
		return
	}

	target, err := h.cfg.DB.User.Get(r.Context(), id)
	if err != nil {
		chix.Error(w, r, err)
		return
	}

	if target.Admin {
		chix.ErrorWithCode(w, r, http.StatusBadRequest, errors.New("cannot ban a system admin"))
		return
	}

	target, err = target.Update().
		SetBanned(true).
		SetBanReason(body.Reason).
		SetBannedByID(self.ID).
		Save(r.Context())
	if err != nil {
		chix.Error(w, r, err)
		return
	}

	chix.JSON(w, r, http.StatusOK, target)
}

func (h *handler) postUnbanUser(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDParam(w, r)
	if !ok {
		return
	}

	target, err := h.cfg.DB.User.UpdateOneID(id).
		SetBanned(false).
		ClearBanReason().
		ClearBannedBy().
		Save(r.Context())
	if err != nil {
		chix.Error(w, r, err)
		return
	}

	chix.JSON(w, r, http.StatusOK, target)
}
