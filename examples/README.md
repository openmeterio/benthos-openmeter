# Examples

The following examples demonstrate how to ingest events from various sources into [OpenMeter](https://openmeter.io) using [Benthos](https://benthos.dev).

The examples use the custom Benthos distribution in this repository.

- [HTTP server](http-server/) (forwarding events to OpenMeter)
- [Kubernetes Pod execution time](kubernetes-pod-exec-time/)

## Prerequisites

These examples require a running OpenMeter instance. If you don't have one, you can [sign up for a free trial](https://openmeter.cloud/sign-up).

If you are using OpenMeter Cloud, [grab a new API token](https://openmeter.cloud/ingest).

### Create a meter

In order to see data in OpenMeter, you need to create a meter first.

In OpenMeter Cloud, go to the [Meters](https://openmeter.cloud/meters) page and click the **Create meter** button in the right upper corner.
Fill in the details of the meter as instructed by the specific example and click **Create**.

> [!TIP]
> You can start ingesting events without creating a meter first, but you won't be able to query data.
> You can inspect the ingested events in the [Event debugger](https://openmeter.cloud/ingest/debug).

In a self-hosted OpenMeter instance you can create a meter in the configuration file:

```yaml
# ...

meters:
  - slug: api_calls
    eventType: api-calls
    aggregation: SUM
    valueProperty: $.duration_ms
    groupBy:
      method: $.method
      path: $.path
```
