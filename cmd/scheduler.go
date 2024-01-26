package main

import (
	"os"

	"example.com/custom-scheduler/plugins/nodenumber"
	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/logs/json/register" // for JSON log format registration
	_ "k8s.io/component-base/metrics/prometheus/clientgo"
	_ "k8s.io/component-base/metrics/prometheus/version" // for version metric registration
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(nodenumber.Name, nodenumber.New),
	)

	code := cli.Run(command)
	os.Exit(code)
}
