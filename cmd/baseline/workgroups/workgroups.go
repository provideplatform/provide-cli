package workgroups

import (
	"github.com/provideplatform/provide-cli/cmd/baseline/participants"
	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var WorkgroupsCmd = &cobra.Command{
	Use:   "workgroups",
	Short: "Interact with a baseline workgroups",
	Long:  `Create, manage and interact with workgroups via the baseline protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

func init() {
	WorkgroupsCmd.AddCommand(initBaselineWorkgroupCmd)
	WorkgroupsCmd.AddCommand(joinBaselineWorkgroupCmd)
	WorkgroupsCmd.AddCommand(listBaselineWorkgroupsCmd)
	WorkgroupsCmd.AddCommand(participants.ParticipantsCmd)
	WorkgroupsCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
	WorkgroupsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
