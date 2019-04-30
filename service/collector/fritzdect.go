package collector

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schjan/fritzdect-exporter/client"
	"sync"
)

type FritzDectConfig struct {
	// Dependencies.
	Logger micrologger.Logger
	Client client.Client
}

type FritzDectCollector struct {
	client client.Client

	// Internals.
	bootOnce sync.Once

	currentTempMetric *prometheus.Desc
	desiredTempMetric *prometheus.Desc
	batteryMetric     *prometheus.Desc
}

func NewFritzDectCollector(config FritzDectConfig) (*FritzDectCollector, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Client == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Client must not be empty", config)
	}

	collector := FritzDectCollector{
		client: config.Client,

		currentTempMetric: prometheus.NewDesc("fritz_dect_temp_current",
			"Current room temperature measured by thermostat",
			[]string{"name"}, nil),
		desiredTempMetric: prometheus.NewDesc("fritz_dect_temp_desired",
			"Desired temperature of thermostat",
			[]string{"name"}, nil),
		batteryMetric: prometheus.NewDesc("fritz_dect_batterycharge",
			"Remaining battery charge of device",
			[]string{"name"}, nil),
	}

	return &collector, nil
}

func (c *FritzDectCollector) Describe(ch chan<- *prometheus.Desc) error {
	ch <- c.desiredTempMetric
	ch <- c.currentTempMetric
	ch <- c.batteryMetric

	return nil
}

func (c *FritzDectCollector) Collect(ch chan<- prometheus.Metric) error {
	i, err := c.client.GetDeviceListInfos()
	if err != nil {
		if client.IsUnauthenticated(err) {
			err = c.client.Login()
			if err != nil {
				return microerror.Maskf(err, "was unauthenticated, tried to login, that failed")
			}
			i, err = c.client.GetDeviceListInfos()
			if err != nil {
				return microerror.Mask(err)
			}
		} else {
			return microerror.Mask(err)
		}
	}

	for _, d := range i.Device {
		if d.Temperature != nil {
			ch <- prometheus.MustNewConstMetric(c.currentTempMetric, prometheus.GaugeValue, client.TemperatureToFloat(d.Temperature.Celsius+d.Temperature.Offset), d.Name)
			ch <- prometheus.MustNewConstMetric(c.desiredTempMetric, prometheus.GaugeValue, client.WeirdTemperatureToFloat(d.Hkr.Tsoll), d.Name)
			ch <- prometheus.MustNewConstMetric(c.batteryMetric, prometheus.GaugeValue, float64(d.Hkr.Battery), d.Name)
		}
	}

	return nil
}
