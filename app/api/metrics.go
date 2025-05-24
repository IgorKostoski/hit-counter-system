package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	hitsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "hitcounter_api",
			Name:      "hits_processed_total",
			Help:      "Total number of hit processing attempts.",
		},
		[]string{"status"},
	)
)
