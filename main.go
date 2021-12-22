package main

import (
	"os"

	"github.com/sensu/sensu-go/types"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-plugin-sdk/sensu/metric"
)

// Config represents the handler plugin config.
type Config struct {
	sensu.PluginConfig
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "dump-to-prom",
			Short:    "dumps sensu event check metrics to prometheus exposition",
			Keyspace: "sensu.io/plugins/dump-to-prom/config",
		},
	}

	options = []*sensu.PluginConfigOption{}
)

func main() {
	handler := sensu.NewGoHandler(&plugin.PluginConfig, options, checkArgs, executeHandler)
	handler.Execute()
}

func checkArgs(_ *types.Event) error {
	return nil
}

func executeHandler(event *types.Event) error {
	if !event.HasMetrics() {
		return nil
	}

	f, err := os.OpenFile("/tmp/promout.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	metric.Points(event.Metrics.Points).ToProm(f)
	return nil
}
