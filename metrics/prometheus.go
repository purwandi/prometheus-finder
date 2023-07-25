package metrics

import "github.com/prometheus/client_golang/prometheus"

func DumpFile(path string, collector []prometheus.Collector) error {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector...)

	return prometheus.WriteToTextfile(path, registry)
}
