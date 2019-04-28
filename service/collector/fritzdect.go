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
			[]string{"room"}, nil),
		desiredTempMetric: prometheus.NewDesc("fritz_dect_temp_desired",
			"Desired temperature of thermostat",
			[]string{"room"}, nil),
	}

	return &collector, nil
}

func (c *FritzDectCollector) Describe(ch chan<- *prometheus.Desc) error {
	ch <- c.desiredTempMetric
	ch <- c.currentTempMetric

	return nil
}

func (c *FritzDectCollector) Collect(ch chan<- prometheus.Metric) error {
	ch <- prometheus.MustNewConstMetric(c.desiredTempMetric, prometheus.GaugeValue, 21.2, "dorm")
	ch <- prometheus.MustNewConstMetric(c.desiredTempMetric, prometheus.GaugeValue, 15.6, "kitchen")

	ch <- prometheus.MustNewConstMetric(c.currentTempMetric, prometheus.GaugeValue, 25, "dorm")
	ch <- prometheus.MustNewConstMetric(c.currentTempMetric, prometheus.GaugeValue, 18, "kitchen")

	return nil
}
