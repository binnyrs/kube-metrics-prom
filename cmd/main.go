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

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"

	"github.com/bmcstdio/kube-metrics-prom/pkg/collector"
)

var (
	fs *flag.FlagSet

	bindHost   string
	bindPort   int
	kubeconfig string
)

func init() {
	fs = flag.NewFlagSet("", flag.ExitOnError)

	fs.StringVar(&bindHost, "bind-host", "0.0.0.0", "host to bind to")
	fs.IntVar(&bindPort, "bind-port", 9100, "port to bind to")
	fs.StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig (if running outside the cluster)")

	// parse provided flags
	fs.Parse(os.Args[1:])
	// https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})
}

func main() {
	// say hi
	log.Info("kube-metrics-prom is starting")
	// read the specified kubeconfig file
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("failed to read kubeconfig: %v", err)
	}
	// create a client for the metrics.k8s.io/v1beta1 api
	metricsClient, err := metricsv1beta1.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("failed to build metrics.k8s.io/v1beta1 client: %v", err)
	}
	// create and register our collector for metrics.k8s.io/v1beta1 metrics
	if err := prometheus.Register(collector.NewMetricsV1Beta1Collector(metricsClient)); err != nil {
		log.Fatalf("failed to register the collector for metrics.k8s.io/v1beta1: %v", err)
	}
	// expose a dummy endpoint under /healthz
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	// expose prometheus metrics under /metrics
	http.Handle("/metrics", promhttp.Handler())
	// create and run the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", bindHost, bindPort), nil))
}
