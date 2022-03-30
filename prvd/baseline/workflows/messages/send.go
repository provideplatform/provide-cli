package messages

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var baselineAPIEndpoint string
var baselineID string
var data string
var id string
var messageType string
var recipients string
var Optional bool

var sendBaselineMessageCmd = &cobra.Command{
	Use:   "send",
	Short: "Send baseline message",
	Long:  `Send baseline message in the context of a workflow`,
	Run:   sendMessage,
}

func sendMessage(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepSend)
}

func sendMessageRun(cmd *cobra.Command, args []string) {
	common.AuthorizeApplicationContext()
	common.AuthorizeOrganizationContext(false)

	var payload map[string]interface{}
	err := json.Unmarshal([]byte(data), &payload)
	if err != nil {
		log.Printf("WARNING: failed to send baseline message; failed to parse message data as JSON; %s", err.Error())
		os.Exit(1)
	}

	params := map[string]interface{}{
		"id":      id,
		"payload": payload,
		"type":    messageType,
	}
	if baselineID != "" {
		params["baseline_id"] = baselineID
	}
	if recipients != "" {
		_recipients := make([]*baseline.Participant, 0)
		for _, id := range strings.Split(recipients, ",") {
			orgs, err := ident.ListApplicationOrganizations(common.ApplicationAccessToken, common.ApplicationID, map[string]interface{}{
				"organization_id": id,
			})
			if err != nil {
				log.Printf("WARNING: failed to send message message data as JSON; %s", err.Error())
				os.Exit(1)
			}
			for _, org := range orgs {
				if addr, addrOk := org.Metadata["address"].(string); addrOk {
					_recipients = append(_recipients, &baseline.Participant{
						Address: &addr,
					})
				}

			}
		}
		params["recipients"] = _recipients
	}

	baselinedRecord, err := baseline.SendProtocolMessage(common.OrganizationAccessToken, params)
	if err != nil {
		log.Printf("WARNING: failed to baseline %d-byte payload; %s", len(data), err.Error())
		os.Exit(1)
	}

	log.Printf("baselined record: %v", baselinedRecord.(map[string]interface{})["baseline_id"].(string))
	if common.Verbose {
		raw, _ := json.MarshalIndent(baselinedRecord, "", "  ")
		log.Printf(string(raw))
	}

}

func init() {
	// runBaselineStackCmd.Flags().StringVar(&baselineAPIEndpoint, "baseline-api-endpoint", "", "baseline API endpoint for use by one or more authorized systems of record")

	sendBaselineMessageCmd.Flags().StringVar(&baselineID, "baseline-id", "", "the globally-unique baseline identifier for the record")

	sendBaselineMessageCmd.Flags().StringVar(&data, "data", "", "content of the message")
	// sendBaselineMessageCmd.MarkFlagRequired("data")

	sendBaselineMessageCmd.Flags().StringVar(&id, "id", "", "identifier of the associated payload in the internal system of record")
	// sendBaselineMessageCmd.MarkFlagRequired("id")

	sendBaselineMessageCmd.Flags().StringVar(&messageType, "type", "", "type of the payload to be baselined")
	// sendBaselineMessageCmd.Flags().StringVar(&recipients, "recipients", "", "comma-delimited list of recipient organization ids")

	sendBaselineMessageCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	// sendBaselineMessageCmd.MarkFlagRequired("organization")

	sendBaselineMessageCmd.Flags().StringVar(&common.ApplicationID, "workgroup", "", "workgroup identifier")
	//sendBaselineMessageCmd.MarkFlagRequired("workgroup")
}
