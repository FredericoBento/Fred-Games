package app

import "errors"

var (
	ErrServerRouterNotFound = errors.New("server router could not be found in app")
)
