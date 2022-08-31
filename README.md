[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/jadiunr/check-memory-usage-extended)
![Go Test](https://github.com/jadiunr/check-memory-usage-extended/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/jadiunr/check-memory-usage-extended/workflows/goreleaser/badge.svg)

# Sensu memory and swap usage checks extended

## Table of Contents
- [Overview](#overview)
  - [Checks](#checks)
- [Usage examples](#usage-examples)
  - [check-memory-usage-extended](#check-memory-usage-extended)
  - [check-swap-usage-extended](#check-swap-usage-extended)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The Sensu memory and swap usage checks extended are a [Sensu Check][6] that provide alerting and more detailed metrics for memory and swap usage. Metrics are provided in [nagios_perfdata](https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/perfdata.html) (default) or [influxdb_line](https://docs.influxdata.com/enterprise_influxdb/v1.9/write_protocols/line_protocol_reference/)

### Checks

This collection contains the following checks:

- check-memory-usage-extended - for checking memory usage
- check-swap-usage-extended - for checking swap usage

## Usage examples

### check-memory-usage-extended

```
Check memory usage in more detail

Usage:
  check-memory-usage-extended [flags]
  check-memory-usage-extended [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -c, --critical float   Critical threshold for overall memory usage (default 90)
  -f, --format string    Choose output format 'nagios' or 'influxdb' (default "nagios")
  -h, --help             help for check-memory-usage-extended
  -w, --warning float    Warning threshold for overall memory usage (default 75)

Use "check-memory-usage-extended [command] --help" for more information about a command.
```

### check-swap-usage-extended

```
Check swap usage in more detail

Usage:
  check-swap-usage-extended [flags]
  check-swap-usage-extended [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -c, --critical float   Critical threshold for overall swap usage (default 90)
  -f, --format string    Choose output format 'nagios' or 'influxdb' (default "nagios")
  -h, --help             help for check-swap-usage-extended
  -w, --warning float    Warning threshold for overall swap usage (default 75)

Use "check-swap-usage-extended [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add jadiunr/check-memory-usage-extended
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/jadiunr/check-memory-usage-extended].

### Check definitions

#### check-memory-usage-extended

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: check-memory-usage-extended
  namespace: default
spec:
  command: check-memory-usage-extended --critical 90 --warning 80
  output_metric_format: nagios_perfdata
  output_metric_handlers:
  - influxdb
  subscriptions:
  - system
  runtime_assets:
  - jadiunr/check-memory-usage-extended
```

#### check-swap-usage-extended

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: check-swap-usage-extended
  namespace: default
spec:
  command: check-swap-usage-extended --critical 90 --warning 75
  output_metric_format: nagios_perfdata
  output_metric_handlers:
  - influxdb
  subscriptions:
  - system
  runtime_assets:
  - jadiunr/check-memory-usage-extended
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the check-memory-usage-extended repository:

```
go build ./cmd/check-memory-usage-extended/
go build ./cmd/check-swap-usage-extended/
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/jadiunr/check-memory-usage-extended/blob/master/.github/workflows/release.yml
[5]: https://github.com/jadiunr/check-memory-usage-extended/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[7]: https://github.com/sensu/check-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
