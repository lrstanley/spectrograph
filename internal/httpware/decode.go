// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
)

var globalDecoder = form.NewDecoder()

func DecodeData(w http.ResponseWriter, r *http.Request, v interface{}) (ok bool) {
	var err, rerr error

	if err = r.ParseForm(); err != nil {
		rerr = fmt.Errorf("error parsing %s parameters, invalid request", r.Method)
	} else {
		switch r.Method {
		case http.MethodGet:
			err = globalDecoder.Decode(v, r.Form)
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				dec := json.NewDecoder(r.Body)
				defer r.Body.Close()
				err = dec.Decode(v)
			} else {
				err = globalDecoder.Decode(v, r.PostForm)
			}
		}
		if err != nil {
			rerr = fmt.Errorf("error decoding %s request into required format (%T): validate request parameters", r.Method, v)
		}
	}

	if err != nil {
		_ = Error(w, r, http.StatusBadRequest, rerr)
	}

	return err == nil
}
