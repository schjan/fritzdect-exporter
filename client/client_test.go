package client

import "testing"

func TestLogin(t *testing.T) {
	c := client{
		rootUrl: "http://fritz.box",
	}

	c.login()
}
