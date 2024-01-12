# HTTP server

This example demonstrates how to run an HTTP server that accepts and forwards events to OpenMeter.
It enables the operation of a low-latency ingestion point near your services, which also handles retries, back pressure, and buffering (if necessary).

The nice thing about this solution is that you can use it as a drop-in replacement for existing integrations since it's compatible with the [OpenMeter Ingest API](https://openmeter.io/docs/getting-started/rest-api).
This means you can use our existing SDKs; the only thing you need to change is the API endpoint.
That is, of course, optional: you can use any client library or payload format you prefer and handle the mapping to the [CloudEvents](https://cloudevents.io/) format in a Benthos [`mapping` processor](https://www.benthos.dev/docs/components/processors/mapping/).

## Table of Contents

- [Prerequisites](#prerequisites)
- [Launch the example](#launch-the-example)
- [Cleanup](#cleanup)
- [Advanced configuration](#advanced-configuration)
- [Production use](#production-use)

## Prerequisites

This example uses [Docker](https://docker.com) and [Docker Compose](https://docs.docker.com/compose/), but you are free to run the components in any other way.

Check out this repository if you want to run the example locally:

```shell
git clone https://github.com/openmeterio/benthos-openmeter.git
cd benthos-openmeter/examples/http-server
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
> Read more about creating a meter in the general examples [README](../../README.md#Create-a-meter).

## Launch the example

Launch the event forwarder:

```shell
docker compose up -d
```

Start tailing the logs:

```shell
docker compose logs -f forwarder
```

Open another terminal and send a test event:

```shell
curl -vvv http://localhost:4196/api/v1/events -H "Content-Type: application/cloudevents+json" -d @seed/event.json
```

_Inspect the logs in the other terminal._

Send a batch of events:

```shell
curl -vvv http://localhost:4196/api/v1/events -H "Content-Type: application/cloudevents-batch+json" -d @seed/batch.json
```

_Inspect the logs in the other terminal._

Start the seeder component to continuously send events to the forwarder:

```shell
docker compose --profile seed up -d
```

> [!WARNING]
> The seeder component sends one event per second to the forwarder.

> [!NOTE]
> The seeder sends events to the forwarder from 3 days ago up to the current time.
> For this reason, events may not immediately appear in the event debugger.
>
> You can modify this behavior by editing the [seeder configuration file](seed/config.yaml).

## Checking events

Read more in the general examples [README](../../README.md#Checking-events-in-OpenMeter).

## Cleanup

Stop containers:

```shell
docker compose --profile seed down -v
```

## Advanced configuration

For this example to work, Benthos needs to be run in [streams mode](https://www.benthos.dev/docs/guides/streams_mode/about).

_(The reason is that the strong delivery guarantees currently do not permit the HTTP server to reply before events are forwarded to OpenMeter, which may lead to timeouts.
To learn more, watch [this discussion](https://github.com/benthosdev/benthos/discussions/2324).)_

As a result, this example requires two configuration files:

- [input.yaml](input.yaml) contains the configuration of the HTTP server and validation of incoming events.
- [output.yaml](output.yaml) contains the configuration of the OpenMeter output.

Run Benthos with these configuration files:

```shell
benthos streams input.yaml output.yaml
```

Check out the configuration files and the [Benthos documentation](https://www.benthos.dev/docs/about) for more details.

## Production use

We are actively working on improving the documentation and the examples.
In the meantime, feel free to contact us [in email](https://us10.list-manage.com/contact-form?u=c7d6a96403a0e5e19032ee885&form_id=fe04a7fc4851f8547cfee56763850e95) or [on Discord](https://discord.gg/nYH3ZQ3Xzq).

We are more than happy to help you set up OpenMeter in your production environment.
