package client

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDeviceListInfos(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//assert.Contains(t, req.RequestURI, expected)

		bts, err := ioutil.ReadFile("./deviceListExample.xml")
		if err != nil {
			t.Fatal(err)
		}

		rw.Write(bts)
	}))
	defer server.Close()

	c := client{
		url:  server.URL,
		http: server.Client(),
	}

	info, err := c.GetDeviceListInfos()
	assert.NoError(t, err)
	assert.Equal(t, 220, info.Device[0].Temperature.Celsius)
}
