package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/schjan/fritzdect-exporter/collector"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())

	log.Info("Beginning to serve on :8080")

	c, err := collector.NewFritzDect()
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(c)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
