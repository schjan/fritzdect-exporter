package collector

import "github.com/prometheus/client_golang/prometheus"

type fritzdectCollector struct {
	currentTempMetric *prometheus.Desc
	desiredTempMetric *prometheus.Desc
}

func NewFritzDectCollector() (*fritzdectCollector, error) {
	return &fritzdectCollector{
		currentTempMetric: prometheus.NewDesc("fritz_dect_temp_current",
			"Current room temperature measured by thermostat",
			[]string{"room"}, nil),
		desiredTempMetric: prometheus.NewDesc("fritz_dect_temp_desired",
			"Desired temperature of thermostat",
			[]string{"room"}, nil),
	}, nil
}

func (c *fritzdectCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *fritzdectCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.desiredTempMetric, prometheus.GaugeValue, 21.2, "dorm")
	ch <- prometheus.MustNewConstMetric(c.desiredTempMetric, prometheus.GaugeValue, 15.6, "kitchen")

	ch <- prometheus.MustNewConstMetric(c.currentTempMetric, prometheus.GaugeValue, 25, "dorm")
	ch <- prometheus.MustNewConstMetric(c.currentTempMetric, prometheus.GaugeValue, 18, "kitchen")
}
