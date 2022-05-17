/*
 * Copyright 2017-2022 Provide Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package stack

import (
	"log"
	"os"
	"sync"

	"github.com/docker/docker/client"
	"github.com/provideplatform/provide-cli/prvd/common"

	"github.com/spf13/cobra"
)

var logsBaselineStackCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print baseline stack logs",
	Long:  `Print the logs from each container in a local baseline stack instance`,
	Run:   stackLogs,
}

func stackLogs(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepLogs)
}

func stackLogsRun(cmd *cobra.Command, args []string) {
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	common.LogContainers(docker, &wg, name)
	wg.Wait()
}

func init() {
	logsBaselineStackCmd.Flags().StringVar(&name, "name", "baseline-local", "name of the baseline stack instance")
}
