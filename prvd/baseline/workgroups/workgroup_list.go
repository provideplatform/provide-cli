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

package workgroups

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	ident "github.com/provideplatform/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var listBaselineWorkgroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List baseline workgroups",
	Long:  `List all available baseline workgroups`,
	Run:   listWorkgroups,
}

func listWorkgroups(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listWorkgroupsRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	applications, err := ident.ListApplications(token, map[string]interface{}{
		"type": "baseline",
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	})
	if err != nil {
		log.Printf("failed to retrieve baseline workgroups; %s", err.Error())
		os.Exit(1)
	}
	for i := range applications {
		workgroup := applications[i]
		result := fmt.Sprintf("%s\t%s\n", workgroup.ID.String(), *workgroup.Name)
		fmt.Print(result)
	}
}

func init() {
	listBaselineWorkgroupsCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	listBaselineWorkgroupsCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of baseline workgroups to retrieve per page")
	listBaselineWorkgroupsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
