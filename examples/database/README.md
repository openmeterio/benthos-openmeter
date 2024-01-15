# Database

This example demonstrates reading data from a database, transforming it to [CloudEvents](https://cloudevents.io/) and sending to [OpenMeter](https://openmeter.io/).

This is a rather common use case when a system already collects some sort of data or log and you want to send that as usage data to OpenMeter for further processing.

In this example, we will read data from a fake chat service's database (automatically seeded with random data) that charges customers based on the number of messages and/or message length.
Benthos will read messages from a message log table and send the calculated usage to OpenMeter.

The example also demonstrates that certain business logic can also be implemented during the transformation (for example: users on the enterprise plan do not get charged for message length).

Databases featured in this example:

- Postgres
- [Clickhouse](https://clickhouse.com/)

> [!TIP]
> Check out the supported database drivers in the [Benthos documentation](https://www.benthos.dev/docs/components/inputs/sql_select#drivers).

## Table of Contents <!-- omit from toc -->

- [Prerequisites](#prerequisites)
- [Launch the example](#launch-the-example)
- [Checking events](#checking-events)
- [Cleanup](#cleanup)
- [Advanced configuration](#advanced-configuration)
- [Production use](#production-use)

## Prerequisites

This example uses [Docker](https://docker.com) and [Docker Compose](https://docs.docker.com/compose/), but you are free to run the components in any other way.

Check out this repository if you want to run the example locally:

```shell
git clone https://github.com/openmeterio/benthos-openmeter.git
cd benthos-openmeter/examples/database
```

Create a new `.env` file and add the details of your OpenMeter instance:

```shell
cp .env.dist .env
# edit .env and fill in the details
```

> [!TIP]
> Tweak other options in the `.env` file to change the behavior of the example.

Create the following meters in OpenMeter with the following details:

| Description              | Event type     | Aggregation | Value property              | Group by (optional) |
| ------------------------ | -------------- | ----------- | --------------------------- | ------------------- |
| All message lenghts      | `chat-message` | `SUM`       | `$.message_length`          | - `plan`: `$.plan`  |
| Billable message lengths | `chat-message` | `SUM`       | `$.message_length_billable` | - `plan`: `$.plan`  |
| Message count            | `chat-message` | `COUNT`     | -                           | - `plan`: `$.plan`  |

> [!TIP]
> Read more about creating a meter in the general examples [README](../README.md#Create-a-meter).

## Launch the example

Decide which database you want to use:

```shell
export COMPOSE_PROFILES=SELECTED_DATABASE
```

Available profiles:

- `postgres`
- `clickhouse`

Launch the example (database, event collector and seeder):

```shell
docker compose up -d
```


## Checking events

Read more in the general examples [README](../README.md#Checking-events-in-OpenMeter).

## Cleanup

Stop containers:

```shell
docker compose down -v
```

## Advanced configuration

Check out the configuration files and the [Benthos documentation](https://www.benthos.dev/docs/about) for more details.

## Production use

We are actively working on improving the documentation and the examples.
In the meantime, feel free to contact us [in email](https://us10.list-manage.com/contact-form?u=c7d6a96403a0e5e19032ee885&form_id=fe04a7fc4851f8547cfee56763850e95) or [on Discord](https://discord.gg/nYH3ZQ3Xzq).

We are more than happy to help you set up OpenMeter in your production environment.
