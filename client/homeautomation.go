package client

import (
	"encoding/xml"
	"github.com/giantswarm/microerror"
	"math"
	"net/http"
)

const (
	webservicesUrl    = "webservices"
	homeautoSwitchUrl = webservicesUrl + "/homeautoswitch.lua"

	getDeviceListInfosCmd = "getdevicelistinfos"
	switchCmd             = "switchcmd"
)

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

func (c *client) GetDeviceListInfos() (*DeviceListInfo, error) {
	resp, err := c.r().
		SetQueryParam(switchCmd, getDeviceListInfosCmd).
		Get(homeautoSwitchUrl)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if resp.StatusCode() == http.StatusForbidden {
		return nil, microerror.Mask(unauthenticatedError)
	}

	var deviceList DeviceListInfo
	err = xml.Unmarshal(resp.Body(), &deviceList)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &deviceList, nil
}

func TemperatureToFloat(temp int) float64 {
	return float64(temp) * 0.1
}

func WeirdTemperatureToFloat(temp int) float64 {
	// temp <= 16  -> <= 8°C
	// temp*0.5    -> °C
	// temp >= 56  -> >= 28°C
	// temp == 253 -> OFF
	// temp == 254 -> ON

	if temp == 253 {
		return math.Inf(0)
	}
	if temp == 254 {
		return math.Inf(1)
	}

	if temp <= 16 {
		return 8
	}
	if temp >= 56 {
		return 28
	}

	return float64(temp) * 0.5
}
