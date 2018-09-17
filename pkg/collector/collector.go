/*
Copyright 2018 The prom-metrics-server Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

var (
	// nodeLabels is the set of labels placed in node metrics
	nodeLabels = []string{"name"}
	// podLabels is the set of labels placed in pod metrics
	podLabels = []string{"container", "name", "namespace"}
)

// metricsv1beta1Collector is the struct that backs the metrics.k8s.io/v1beta1 collector.
type metricsv1beta1Collector struct {
	metricsClient *metricsv1beta1.MetricsV1beta1Client
	nodeCPUDesc   *prometheus.Desc
	nodeMemDesc   *prometheus.Desc
	podCPUDesc    *prometheus.Desc
	podMemDesc    *prometheus.Desc
}

// NewMetricsV1Beta1Collector returns a Prometheus collector that exposes CPU and memory metrics for nodes and pods.
// Metrics are sourced from the metrics.k8s.io/v1beta1 API at the time of collection.
func NewMetricsV1Beta1Collector(metricsClient *metricsv1beta1.MetricsV1beta1Client) prometheus.Collector {
	return &metricsv1beta1Collector{
		metricsClient: metricsClient,
		nodeCPUDesc:   prometheus.NewDesc("metricsv1beta1_node_cpu", "node cpu usage metrics", nodeLabels, nil),
		nodeMemDesc:   prometheus.NewDesc("metricsv1beta1_node_mem", "node memory usage metrics", nodeLabels, nil),
		podCPUDesc:    prometheus.NewDesc("metricsv1beta1_pod_cpu", "pod cpu usage metrics", podLabels, nil),
		podMemDesc:    prometheus.NewDesc("metricsv1beta1_pod_mem", "pod memory usage metrics", podLabels, nil),
	}
}

func (mc *metricsv1beta1Collector) Collect(ch chan<- prometheus.Metric) {
	// collect node metrics
	n, err := mc.metricsClient.NodeMetricses().List(v1.ListOptions{})
	if err != nil {
		ch <- prometheus.NewInvalidMetric(mc.nodeCPUDesc, err)
		ch <- prometheus.NewInvalidMetric(mc.nodeMemDesc, err)
	} else {
		for _, node := range n.Items {
			ch <- prometheus.MustNewConstMetric(
				mc.nodeCPUDesc,
				prometheus.GaugeValue,
				float64(node.Usage.Cpu().MilliValue()),
				node.Name)
			ch <- prometheus.MustNewConstMetric(
				mc.nodeMemDesc,
				prometheus.GaugeValue,
				float64(node.Usage.Memory().MilliValue()),
				node.Name)
		}
	}
	// collect pod metrics across all namespaces
	p, err := mc.metricsClient.PodMetricses(v1.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		ch <- prometheus.NewInvalidMetric(mc.podCPUDesc, err)
		ch <- prometheus.NewInvalidMetric(mc.podMemDesc, err)
	} else {
		for _, pod := range p.Items {
			for _, container := range pod.Containers {
				ch <- prometheus.MustNewConstMetric(
					mc.podCPUDesc,
					prometheus.GaugeValue,
					float64(container.Usage.Cpu().MilliValue()),
					container.Name,
					pod.Name,
					pod.Namespace)
				ch <- prometheus.MustNewConstMetric(
					mc.podMemDesc,
					prometheus.GaugeValue,
					float64(container.Usage.Memory().MilliValue()),
					container.Name,
					pod.Name,
					pod.Namespace)
			}
		}
	}
}

func (mc *metricsv1beta1Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- mc.podCPUDesc
	ch <- mc.podMemDesc
}
