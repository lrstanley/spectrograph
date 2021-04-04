// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import "gopkg.in/go-playground/validator.v9"

var validate = validator.New()

// Omit can be used to ignore fields on a given struct, when wrapping the struct.
// For example, given the following structs:
// type User struct {
//     Email    string `json:"email"`
//     Password string `json:"password"`
//     // many more fields…
// }
// type PublicUser struct {
//     *User
//     Token    string `json:"token"`
//     Password models.Omit   `json:"password,omitempty"`
// }
// Using PublicUser{User}, would ensure Password isn't exposed.
type Omit *struct{}
