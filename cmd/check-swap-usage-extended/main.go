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

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-swap-usage-extended",
			Short:    "Check swap usage in more detail",
			Keyspace: "sensu.io/plugins/check-memory-usage-extended/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[float64]{
			Path:      "critical",
			Argument:  "critical",
			Shorthand: "c",
			Default:   float64(90),
			Usage:     "Critical threshold for overall swap usage",
			Value:     &plugin.Critical,
		},
		&sensu.PluginConfigOption[float64]{
			Path: "warning",
			Argument: "warning",
			Shorthand: "w",
			Default: float64(75),
			Usage: "Warning threshold for overall swap usage",
			Value: &plugin.Warning,
		},
		&sensu.PluginConfigOption[string]{
			Path: "format",
			Argument: "format",
			Shorthand: "f",
			Default: "nagios",
			Usage: "Choose output format 'nagios' or 'influxdb'",
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
	vmStat, err := mem.SwapMemory()
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("failed to get swap statistics: %v", err)
	}

	v := reflect.ValueOf(vmStat)
	v = reflect.Indirect(v)

	var output string
	var status int

	switch plugin.Format {
	case "nagios":
		output, status = makeNagiosPerfData(&v)
	case "influxdb":
		output, status = makeInfluxDBLines(&v)
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
		output = append(output, fmt.Sprintf("%s Critical: %.2f%% swap usage | ", plugin.PluginConfig.Name, vmStat.FieldByName("UsedPercent").Float()))
		status = sensu.CheckStateCritical
	} else if vmStat.FieldByName("UsedPercent").Float() > plugin.Warning {
		output = append(output, fmt.Sprintf("%s Warning: %.2f%% swap usage | ", plugin.PluginConfig.Name, vmStat.FieldByName("UsedPercent").Float()))
		status = sensu.CheckStateWarning
	} else {
		output = append(output, fmt.Sprintf("%s OK: %.2f%% swap usage | ", plugin.PluginConfig.Name, vmStat.FieldByName("UsedPercent").Float()))
	}

	for i := 0; i< vmStat.NumField(); i++ {
		measurement := strcase.ToSnake(typeOfVMStat.Field(i).Name)
		perfData = append(perfData, fmt.Sprintf("swap_%s=%v", measurement, vmStat.Field(i).Interface()))
	}

	output = append(output, strings.Join(perfData, ", "))

	return strings.Join(output, ""), status
}

func makeInfluxDBLines(vmStat *reflect.Value) (string, int) {
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
		output = append(output, fmt.Sprintf("swap_%s %v", measurement, vmStat.Field(i).Interface()))
	}

	return strings.Join(output, "\n"), status
}
