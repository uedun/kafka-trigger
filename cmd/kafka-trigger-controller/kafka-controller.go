/*
Copyright (c) 2016-2017 Bitnami

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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	kubelessutils "github.com/kubeless/kubeless/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/kubeless/kafka-trigger/pkg/controller"
	kafkautils "github.com/kubeless/kafka-trigger/pkg/utils"
	"github.com/kubeless/kafka-trigger/pkg/version"
)

const (
	globalUsage = `Kafka trigger controller adds support for Kafka topic as event source to Kubeless functions`
)

var rootCmd = &cobra.Command{
	Use:   "kafka-controller",
	Short: "Kafka controller",
	Long:  globalUsage,
	Run: func(cmd *cobra.Command, args []string) {

		kubelessClient, err := kubelessutils.GetFunctionClientInCluster()
		if err != nil {
			logrus.Fatalf("Cannot get kubeless CR API client: %v", err)
		}

		kafkaClient, err := kafkautils.GetTriggerClientInCluster()
		if err != nil {
			logrus.Fatalf("Cannot get kafka trigger CR API client: %v", err)
		}

		kafkaTriggerCfg := controller.KafkaTriggerConfig{
			KubeCli:        kubelessutils.GetClient(),
			TriggerClient:  kafkaClient,
			KubelessClient: kubelessClient,
		}

		kafkaTriggerController := controller.NewKafkaTriggerController(kafkaTriggerCfg)

		stopCh := make(chan struct{})
		defer close(stopCh)

		go kafkaTriggerController.Run(stopCh)

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGTERM)
		signal.Notify(sigterm, syscall.SIGINT)
		<-sigterm
	},
}

func main() {
	logrus.Infof("Running Kafka controller version: %v", version.Version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
