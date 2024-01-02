# Benthos plugins for OpenMeter

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/openmeterio/benthos-openmeter/ci.yaml?style=flat-square)](https://github.com/openmeterio/benthos-openmeter/actions/workflows/ci.yaml)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/openmeterio/benthos-openmeter)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.20-61CFDD.svg?style=flat-square)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/openmeterio/benthos-openmeter/badge?style=flat-square)](https://api.securityscorecards.dev/projects/github.com/openmeterio/benthos-openmeter)
[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

**A set of plugins and tools for [Benthos](https://www.benthos.dev/) to ingest events into [OpenMeter](https://openmeter.io/).**

## Usage

TODO

## Development

**For an optimal developer experience, it is recommended to install [Nix](https://nixos.org/download.html) and [direnv](https://direnv.net/docs/installation.html).**

_Optional:_ Create a `.env.local` file with the following contents:

```shell
OPENMETER_URL=https://your.openmeter.cloud
OPENMETER_TOKEN=YOUR_TOKEN
```

Run Benthos:

```shell
go run . -c test.yaml
```

## License

The project is licensed under the [Apache 2.0 License](LICENSE).
