# kube-metrics-prom

An exporter of `metrics.k8s.io/v1beta1` metrics in Prometheus format.

## What?

`kube-metrics-prom` exports metrics from the [core metrics pipeline](https://kubernetes.io/docs/tasks/debug-application-cluster/core-metrics-pipeline/) in Prometheus format.

```
# HELP metricsv1beta1_pod_cpu pod cpu usage metrics
# TYPE metricsv1beta1_pod_cpu gauge
metricsv1beta1_pod_cpu{container="coredns",name="coredns-df9bf6587-nrvs9",namespace="kube-system"} 4
metricsv1beta1_pod_cpu{container="kubernetes-dashboard",name="kubernetes-dashboard-7d94764485-j2wbp",namespace="kube-system"} 1
metricsv1beta1_pod_cpu{container="metrics-server",name="metrics-server-6cf7ffcb64-mv5t7",namespace="kube-system"} 4
# HELP metricsv1beta1_pod_mem pod memory usage metrics
# TYPE metricsv1beta1_pod_mem gauge
metricsv1beta1_pod_mem{container="coredns",name="coredns-df9bf6587-nrvs9",namespace="kube-system"} 1.7719296e+10
metricsv1beta1_pod_mem{container="kubernetes-dashboard",name="kubernetes-dashboard-7d94764485-j2wbp",namespace="kube-system"} 1.4917632e+10
metricsv1beta1_pod_mem{container="metrics-server",name="metrics-server-6cf7ffcb64-mv5t7",namespace="kube-system"} 2.3732224e+10
```

## Why?

For learning purposes! 👨‍🏫

Also, there seems to be [some demand](https://github.com/kubernetes-incubator/metrics-server/issues/7) for scraping `metrics-server` using Prometheus.

## How?

It works by directly querying the `metrics.k8s.io/v1beta1` API whenever the `/metrics` endpoint is scraped. 

## Where?

In your local workstation by running

```
$ docker run \
    -p 9100:9100 \
    -v $HOME/.kube/config:/kubeconfig.yaml:ro \
    quay.io/bmcstdio/kube-metrics-prom:v0.1.0 \
    --kubeconfig /kubeconfig.yaml
```

or inside a Kubernetes cluster by running

```
$ kubectl create -f ./examples/kube-metrics-prom.yml
```

## License

Copyright 2018 bmcstdio

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
