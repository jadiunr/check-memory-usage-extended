package main

import (
	"fmt"
	"reflect"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Config struct {
	sensu.PluginConfig
	Critical float64
	Warning float64
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-memory-usage-extended",
			Short:    "Check memory usage in more detail",
			Keyspace: "sensu.io/plugins/check-memory-usage-extended/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[float64]{
			Path:      "critical",
			Argument:  "critical",
			Shorthand: "c",
			Default:   float64(90),
			Usage:     "Critical threshold for overall memory usage",
			Value:     &plugin.Critical,
		},
		&sensu.PluginConfigOption[float64]{
			Path: "warning",
			Argument: "warning",
			Shorthand: "w",
			Default: float64(75),
			Usage: "Warning threshold for overall memory usage",
			Value: &plugin.Warning,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("failed to get virtual memory statistics: %v", err)
	}

	v := reflect.ValueOf(vmStat)
	v = reflect.Indirect(v)

	typeOfVMStat := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("Field: %s\tValue: %v\n", typeOfVMStat.Field(i).Name, v.Field(i).Interface())
	}

	return sensu.CheckStateOK, nil
}
