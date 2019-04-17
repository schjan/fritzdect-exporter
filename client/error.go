package client

import (
	"github.com/giantswarm/microerror"
)

var unauthenticatedError = microerror.New("unauthenticated")

var invalidConfigError = microerror.New("invalid config")

func IsUnauthenticated(err error) bool {
	return err == unauthenticatedError
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}
