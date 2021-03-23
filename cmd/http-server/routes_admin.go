// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"github.com/go-chi/chi"
)

func registerAdminRoutes(r chi.Router) {
	// r.Get("/api/v1/admin/providers", func(w http.ResponseWriter, r *http.Request) {
	// 	resp := pt.M{"providers": []string{}}

	// 	providers, err := svcProviders.List(r.Context())
	// 	if err != nil {
	// 		if models.IsNotFound(err) {
	// 			pt.JSON(w, r, resp)
	// 			return
	// 		}
	// 		helpers.HTTPError(w, r, http.StatusInternalServerError, err)
	// 		return
	// 	}

	// 	resp["providers"] = providers
	// 	pt.JSON(w, r, resp)
	// })
}
