package main

import (
	"os"
	"runtime"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/purwandi/prometheus-finder/config"
	"github.com/purwandi/prometheus-finder/finder"
	"github.com/purwandi/prometheus-finder/metrics"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

func main() {
	var wg sync.WaitGroup
	var collector []prometheus.Collector

	cfgPath, err := config.ParseFlags()
	if err != nil {
		logrus.Fatal(err)
	}

	configs, err := config.NewConfig(cfgPath)
	if err != nil {
		logrus.Fatal(configs)
	}

	if ok, err := config.IsOk(*configs); !ok {
		logrus.Fatal(err)
	}

	collector = append(collector, prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "runtime",
			Name:      "goroutines_count",
			Help:      "Number of goroutines that currently exist.",
		},
		func() float64 { return float64(runtime.NumGoroutine()) },
	))

	for _, item := range *configs {
		wg.Add(1)

		if item.Query == "" {
			continue
		}

		go func(cfg config.Config) {
			defer wg.Done()
			collector = append(collector, finder.Observe(cfg))
		}(item)
	}

	wg.Wait()

	metrics.DumpFile("./cron_exporter.prom", collector)
	defer os.Exit(0)
}
