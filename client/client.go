package client

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Client interface {
	GetDesiredTemperature(ain string) (float32, error)
	GetCurrentTemperature(ain string) (float32, error)
	GetComfortTemperature(ain string) (float32, error)
	GetSavingTemperature(ain string) (float32, error)
}

type client struct {
	http    http.Client
	rootUrl string
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

func (c *client) login() error {
	c.getSid()

	return nil
}

func (c *client) getSid() error {
	resp, err := http.Get(fmt.Sprintf("%s/login_sid.lua", c.rootUrl))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resXml := xml.NewDecoder(resp.Body)

	for {
		t, err := resXml.Token()
		if err != nil {

		}
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "SID" {

			}
		}
	}

	return nil
}
