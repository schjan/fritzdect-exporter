package client

import (
	"net/http"
)

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
	http               *http.Client
	rootUrl            string
	username, password string
}

func New() (*client, error) {
	return nil, nil
}

func (c *client) GetDesiredTemperature(ain string) (float32, error) {
	return 0, nil
}

func (c *client) GetCurrentTemperature(ain string) (float32, error) {
	return 0, nil
}

func (c *client) GetComfortTemperature(ain string) (float32, error) {
	return 0, nil
}

func (c *client) GetSavingTemperature(ain string) (float32, error) {
	return 0, nil
}
