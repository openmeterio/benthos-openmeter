# Generate

This example is based on [our earlier blog post](https://openmeter.io/blog/testing-stream-processing) about testing OpenMeter with Benthos.

The input in this case generates random events (mimicking API calls).

## Table of Contents <!-- omit from toc -->

- [Prerequisites](#prerequisites)
- [Launch the example](#launch-the-example)
- [Checking events](#checking-events)
- [Advanced configuration](#advanced-configuration)
- [Production use](#production-use)

## Prerequisites

Go to the [GitHub Releases](https://github.com/openmeterio/benthos-openmeter/releases/latest) page and download the latest `benthos` binary for your platform.

Check out this repository if you want to run the example locally:

```shell
git clone https://github.com/openmeterio/benthos-openmeter.git
cd benthos-openmeter/examples/generate
```

Create a new `.env` file and add the details of your OpenMeter instance:

```shell
cp .env.dist .env
# edit .env and fill in the details
```

> [!TIP]
> Tweak other options in the `.env` file to change the behavior of the example.

Create a meter in OpenMeter with the following details:

- Event type: `api-calls`
- Aggregation: `SUM`
- Value property: `$.duration_ms`
- Group by (optional):
  - `method`: `$.method`
  - `path`: `$.path`
  - `region`: `$.region`
  - `zone`: `$.zone`

> [!TIP]
> Read more about creating a meter in the general examples [README](../README.md#Create-a-meter).

## Launch the example

Launch the example:

```shell
export OPENMETER_TOKEN=<YOUR TOKEN>
benthos -c config.yaml
```

> [!WARNING]
> By default the example generates 1000 events per second.

## Checking events

Read more in the general examples [README](../README.md#Checking-events-in-OpenMeter).

## Advanced configuration

Check out the configuration files and the [Benthos documentation](https://www.benthos.dev/docs/about) for more details.

## Production use

We are actively working on improving the documentation and the examples.
In the meantime, feel free to contact us [in email](https://us10.list-manage.com/contact-form?u=c7d6a96403a0e5e19032ee885&form_id=fe04a7fc4851f8547cfee56763850e95) or [on Discord](https://discord.gg/nYH3ZQ3Xzq).

We are more than happy to help you set up OpenMeter in your production environment.
