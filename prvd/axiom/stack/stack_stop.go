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

	"github.com/docker/docker/client"
	"github.com/provideplatform/provide-cli/prvd/common"

	"github.com/spf13/cobra"
)

var stopBaselineStackCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the axiom stack",
	Long:  `Stop a local axiom stack instance`,
	Run:   stopStack,
}

func stopStack(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepStop)
}

func runStackStop(cmd *cobra.Command, args []string) {
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	if !prune {
		common.StopContainers(docker, name)
	} else {
		common.PurgeContainers(docker, name, true)
		common.PurgeNetwork(docker, name)
	}

	log.Printf("%s local axiom instance stopped", name)
}

func init() {
	stopBaselineStackCmd.Flags().StringVar(&name, "name", "axiom-local", "name of the axiom stack instance")
	stopBaselineStackCmd.Flags().BoolVar(&prune, "prune", false, "when true, previously-created docker resources are pruned prior to stack initialization")
}
