package client

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"net/http"
)

type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	Username string // Username if user based access is activated
	Password string
	Url      string
}

type Client interface {
	GetDesiredTemperature(ain string) (float32, error)
	GetCurrentTemperature(ain string) (float32, error)
	GetComfortTemperature(ain string) (float32, error)
	GetSavingTemperature(ain string) (float32, error)
}

const (
	unauthenticatedSid = "0000000000000000"


)

type client struct {
	logger micrologger.Logger

	http     *http.Client
	url      string
	username string // username if user based access is activated
	password string

	sid string
}

func New(config Config) (*client, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	// Settings.
	if config.Url == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Url must not be empty", config)
	}
	if config.Password == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Url must not be empty", config)
	}

	s := &client{
		http: http.DefaultClient,

		url:      config.Url,
		username: config.Username,
		password: config.Password,
	}

	return s, nil
}
