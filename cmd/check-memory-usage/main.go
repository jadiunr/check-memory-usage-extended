package main

import (
	"fmt"
	"reflect"
	"strings"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/iancoleman/strcase"
)

type Config struct {
	sensu.PluginConfig
	Critical float64
	Warning float64
	Format string
}

type MetricGroup struct {
	Comment string
	Type string
	Name string
	Value interface{}
}

func (g *MetricGroup) GenerateMetrics() string {
	output := []string{}
	output = append(output, fmt.Sprintf("# HELP mem_%s %s", g.Name, g.Comment))
	output = append(output, fmt.Sprintf("# TYPE mem_%s %s", g.Name, g.Type))
	output = append(output, fmt.Sprintf("mem_%s %v", g.Name, g.Value))
	return strings.Join(output, "\n")
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
		&sensu.PluginConfigOption[string]{
			Path: "format",
			Argument: "format",
			Shorthand: "f",
			Default: "nagios",
			Usage: "Choose output format 'nagios' or 'prometheus'",
			Value: &plugin.Format,
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

	var output string
	var status int

	switch plugin.Format {
	case "nagios":
		output, status = makeNagiosPerfData(&v)
	case "prometheus":
		output, status = makePrometheusMetrics(&v)
	default:
		return sensu.CheckStateCritical, fmt.Errorf("unknown output format: %s", plugin.Format)
	}

	fmt.Println(output)

	return status, nil
}

func makeNagiosPerfData(vmStat *reflect.Value) (string, int) {
	typeOfVMStat := vmStat.Type()
	status := sensu.CheckStateOK
	output := []string{}
	perfData := []string{}

	if vmStat.FieldByName("UsedPercent").Float() > plugin.Critical {
		output = append(output, fmt.Sprintf("%s Critical: %.2f%% memory usage | ", plugin.PluginConfig.Name, vmStat.FieldByName("UsedPercent").Float()))
		status = sensu.CheckStateCritical
	} else if vmStat.FieldByName("UsedPercent").Float() > plugin.Warning {
		output = append(output, fmt.Sprintf("%s Warning: %.2f%% memory usage | ", plugin.PluginConfig.Name, vmStat.FieldByName("UsedPercent").Float()))
		status = sensu.CheckStateWarning
	} else {
		output = append(output, fmt.Sprintf("%s OK: %.2f%% memory usage | ", plugin.PluginConfig.Name, vmStat.FieldByName("UsedPercent").Float()))
	}

	for i := 0; i< vmStat.NumField(); i++ {
		perfData = append(perfData, fmt.Sprintf("mem_%s=%v", strcase.ToSnake(typeOfVMStat.Field(i).Name), vmStat.Field(i).Interface()))
	}

	output = append(output, strings.Join(perfData, ", "))

	return strings.Join(output, ""), status
}

func makePrometheusMetrics(vmStat *reflect.Value) (string, int) {
	typeOfVMStat := vmStat.Type()
	status := sensu.CheckStateOK
	output := []string{}

	if vmStat.FieldByName("UsedPercent").Float() > plugin.Critical {
		status = sensu.CheckStateCritical
	} else if vmStat.FieldByName("UsedPercent").Float() > plugin.Warning {
		status = sensu.CheckStateWarning
	}

	for i := 0; i < vmStat.NumField(); i++ {
		measurement := strcase.ToSnake(typeOfVMStat.Field(i).Name)
		metricGroup := &MetricGroup{
			Name: measurement,
			Type: "untyped",
			Comment: fmt.Sprintf("Statistic %s", typeOfVMStat.Field(i).Name),
			Value: vmStat.Field(i).Interface(),
		}
		output = append(output, metricGroup.GenerateMetrics())
	}

	return strings.Join(output, "\n"), status
}
