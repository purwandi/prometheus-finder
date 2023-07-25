package finder

import (
	"fmt"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/purwandi/prometheus-finder/config"
)

func Observe(config config.Config) prometheus.Collector {
	hostname, _ := os.Hostname()

	if len(config.Labels) == 0 {
		config.Labels = prometheus.Labels{}
	}

	config.Labels["mode"] = string(config.Mode)
	config.Labels["name"] = strings.ToLower(config.Name)
	config.Labels["instance"] = strings.ToLower(hostname)

	opts := prometheus.GaugeOpts{
		Name:        fmt.Sprintf("finder_%s_status", strings.ToLower(config.Name)),
		ConstLabels: config.Labels,
	}

	result, err := ReadFile(config.Path, config.Query, config.Mode)
	if err != nil {
		opts.ConstLabels["at_line"] = ""
		return prometheus.NewGaugeFunc(opts, func() float64 {
			return float64(0)
		})
	}

	opts.ConstLabels["at_line"] = fmt.Sprintf("%d", result.Line)
	return prometheus.NewGaugeFunc(opts, func() float64 {
		return float64(1)
	})
}
