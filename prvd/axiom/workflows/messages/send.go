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

package messages

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/axiom"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var axiomAPIEndpoint string
var axiomID string
var data string
var id string
var messageType string
var recipients string
var Optional bool

var sendBaselineMessageCmd = &cobra.Command{
	Use:   "send",
	Short: "Send axiom message",
	Long:  `Send axiom message in the context of a workflow`,
	Run:   sendMessage,
}

func sendMessage(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepSend)
}

func sendMessageRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}
	if messageType == "" {
		opts := make([]string, 0)
		for k := range items {
			opts = append(opts, k)
		}
		value := common.SelectInput(opts, custodyPromptLabel)
		messageType = items[value]
	}
	if id == "" {
		id = common.FreeInput("ID", "", common.MandatoryValidation)
	}
	if axiomID == "" {
		axiomID = common.FreeInput("Baseline ID", "", common.NoValidation)
	}
	if data == "" {
		data = common.FreeInput("Data", "", common.JSONValidation)
	}

	common.AuthorizeOrganizationContext(true)

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("WARNING: failed to send axiom message; %s", err.Error())
		os.Exit(1)
	}

	var payload map[string]interface{}
	err = json.Unmarshal([]byte(data), &payload)
	if err != nil {
		log.Printf("WARNING: failed to send axiom message; %s", err.Error())
		os.Exit(1)
	}

	params := map[string]interface{}{
		"id":      id,
		"payload": payload,
		"type":    messageType,
	}
	if axiomID != "" {
		params["axiom_id"] = axiomID
	}
	if recipients != "" {
		_recipients := make([]*axiom.Participant, 0)
		for _, id := range strings.Split(recipients, ",") {
			orgs, err := ident.ListApplicationOrganizations(*token.AccessToken, common.ApplicationID, map[string]interface{}{
				"organization_id": id,
			})
			if err != nil {
				log.Printf("WARNING: failed to send message message data as JSON; %s", err.Error())
				os.Exit(1)
			}
			for _, org := range orgs {
				if addr, addrOk := org.Metadata["address"].(string); addrOk {
					_recipients = append(_recipients, &axiom.Participant{
						Address: &addr,
					})
				}

			}
		}
		params["recipients"] = _recipients
	}

	axiomdRecord, err := axiom.SendProtocolMessage(*token.AccessToken, params)
	if err != nil {
		log.Printf("WARNING: failed to axiom %d-byte payload; %s", len(data), err.Error())
		os.Exit(1)
	}

	log.Printf("axiomd record: %v", axiomdRecord.(map[string]interface{})["axiom_id"].(string))
	if common.Verbose {
		raw, _ := json.MarshalIndent(axiomdRecord, "", "  ")
		log.Printf(string(raw))
	}
}

func init() {
	// runBaselineStackCmd.Flags().StringVar(&axiomAPIEndpoint, "axiom-api-endpoint", "", "axiom API endpoint for use by one or more authorized systems of record")

	sendBaselineMessageCmd.Flags().StringVar(&axiomID, "axiom-id", "", "the globally-unique axiom identifier for the record")

	sendBaselineMessageCmd.Flags().StringVar(&data, "data", "", "content of the message")
	// sendBaselineMessageCmd.MarkFlagRequired("data")

	sendBaselineMessageCmd.Flags().StringVar(&id, "id", "", "identifier of the associated payload in the internal system of record")
	// sendBaselineMessageCmd.MarkFlagRequired("id")

	sendBaselineMessageCmd.Flags().StringVar(&messageType, "type", "", "type of the payload to be axiomd")
	// sendBaselineMessageCmd.Flags().StringVar(&recipients, "recipients", "", "comma-delimited list of recipient organization ids")

	sendBaselineMessageCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	// sendBaselineMessageCmd.MarkFlagRequired("organization")

	sendBaselineMessageCmd.Flags().StringVar(&common.ApplicationID, "workgroup", "", "workgroup identifier")
	//sendBaselineMessageCmd.MarkFlagRequired("workgroup")
}
