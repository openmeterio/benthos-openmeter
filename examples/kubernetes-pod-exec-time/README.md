# Kubernetes Pod Execution Time

This example demonstrates metering execution time of Pods running in Kubernetes.

## Prerequisites

Any local (or remote if that's what's available for you) Kubernetes cluster will do.

We will use [kind](https://kind.sigs.k8s.io/) in this example.

Additional tools you are going to need:

- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [helm](https://helm.sh/docs/intro/install/)

## Preparations

Create a new Kubernetes cluster using `kind`:

```shell
kind create cluster
```

> [!TIP]
> Alternatively, set up your `kubectl` context to point to your existing cluster.

Deploy the test Pods to the cluster:

```shell
kubectl apply -f seed/pod.yaml
```

## Deploy the example

Deploy Benthos to your cluster:

```shell
helm install --wait --namespace benthos --create-namespace --set useExample=kubernetes-pod-exec-time --set openmeter.url=<OPENMETER URL> --set openmeter.token=<OPENMETER_TOKEN> benthos-openmeter oci://ghcr.io/openmeterio/helm-charts/benthos-openmeter
```

> [!NOTE]
> If you use OpenMeter Cloud, you can omit the `openmeter.url` parameter.


## Cleanup

Uninstall Benthos from the cluster:

```shell
helm delete --namespace benthos benthos-openmeter
```

Remove the sample Pods from the cluster:

```shell
kubectl delete -f seed/pod.yaml
```

Delete the cluster:

```shell
kind delete cluster
```
