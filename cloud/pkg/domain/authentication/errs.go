// Copyright 2020 Apulis Technology Inc. All rights reserved.

package authentication

import (
	"errors"
)

var (
	ErrAuthTypeNotSupported = errors.New("auth type not supported")
	ErrNoAuthHeader         = errors.New("no auth header")
	ErrAuthenticatorInit    = errors.New("authenticator init fail")
)
