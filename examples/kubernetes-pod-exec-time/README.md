# Kubernetes Pod Execution Time

This example demonstrates metering execution time of Pods running in Kubernetes.

## Prerequisites

Any local (or remote if that's what's available for you) Kubernetes cluster will do.

We will use [kind](https://kind.sigs.k8s.io/) in this example.

Additional tools you are going to need:

- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [helm](https://helm.sh/docs/intro/install/)

## Getting started

Create a new Kubernetes cluster using `kind`:

```shell
kind create cluster
```

> [!TIP]
> Alternatively, set up your `kubectl` context to point to your existing cluster.


## Deploy Benthos

Deploy Benthos to your cluster:

```shell
helm install --wait --namespace benthos --create-namespace benthos-openmeter oci://ghcr.io/openmeterio/helm-charts/benthos-openmeter
```

TODO: add openmeter token
TODO: add kube collector config

Deploy the test Pods to the cluster:

```shell

```

## Cleanup
