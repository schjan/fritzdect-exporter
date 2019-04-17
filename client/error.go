package client

import "github.com/pkg/errors"

var unauthenticatedError = errors.New("unauthenticated")

func IsUnauthenticated(err error) bool {
	return err == unauthenticatedError
}
