package collector

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schjan/fritzdect-exporter/client"
	"sync"
)

type Config struct {
	// Dependencies.
	Logger micrologger.Logger
	Client client.Client
}

type Collector interface {
}

type fritzdect struct {
	client client.Client

	// Internals.
	bootOnce sync.Once

	currentTempMetric *prometheus.Desc
	desiredTempMetric *prometheus.Desc
}

func New(config Config) (*fritzdect, error) {
	collector := &fritzdect{
		currentTempMetric: prometheus.NewDesc("fritz_dect_temp_current",
			"Current room temperature measured by thermostat",
			[]string{"room"}, nil),
		desiredTempMetric: prometheus.NewDesc("fritz_dect_temp_desired",
			"Desired temperature of thermostat",
			[]string{"room"}, nil),
	}

	err := prometheus.Register(collector)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return collector, nil
}

func (c *fritzdect) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *fritzdect) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.desiredTempMetric, prometheus.GaugeValue, 21.2, "dorm")
	ch <- prometheus.MustNewConstMetric(c.desiredTempMetric, prometheus.GaugeValue, 15.6, "kitchen")

	ch <- prometheus.MustNewConstMetric(c.currentTempMetric, prometheus.GaugeValue, 25, "dorm")
	ch <- prometheus.MustNewConstMetric(c.currentTempMetric, prometheus.GaugeValue, 18, "kitchen")
}
