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

package domain_models

import (
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

var DomainModelsCmd = &cobra.Command{
	Use:   "domain-models",
	Short: "Interact with axiom domain models",
	Long:  `Create, manage and interact with domain models via the axiom protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

func init() {
	DomainModelsCmd.AddCommand(listBaselineDomainModelsCmd)
	DomainModelsCmd.AddCommand(initBaselineDomainModelCmd)
	//  SubjectAccountsCmd.AddCommand(subjectAccountDetailsCmd)
	DomainModelsCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
	DomainModelsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
