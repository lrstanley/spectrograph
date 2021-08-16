// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"errors"

	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

// Omit can be used to ignore fields on a given struct, when wrapping the struct.
// For example, given the following structs:
// type User struct {
//     Email    string `json:"email"`
//     Password string `json:"password"`
//     // many more fields...
// }
// type PublicUser struct {
//     *User
//     Token    string `json:"token"`
//     Password models.Omit   `json:"password,omitempty"`
// }
// Using PublicUser{User}, would ensure Password isn't exposed.
type Omit *struct{}

// Validator can be used by a struct to implement it's own validation function,
// rather than using the default struct validator.
type Validator interface {
	Validate() error
}

func Validate(v interface{}) (err error) {
	if custom, ok := v.(Validator); ok {
		err = custom.Validate()

		if err != errUseBuiltinValidator {
			return &ErrValidationFailed{Err: err}
		}
		// Otherwise continue with the default validator.
	}

	err = validate.Struct(v)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}
		return &ErrValidationFailed{Err: err}
	}
	return nil
}

// useBuiltinValidator can be used to tell the validator that it did it's own
// checks, but still wants the standard validator checks to pass.
var errUseBuiltinValidator = errors.New("use builtin")
