package participants

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ParticipantsCmd = &cobra.Command{
	Use:   "participants",
	Short: "Interact with participants in a baseline workgroup",
	Long:  `Invite, manage and interact with workgroup participants via the baseline protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("participants unimplemented")
	},
}

func init() {
	ParticipantsCmd.AddCommand(inviteBaselineWorkgroupParticipantCmd)
	ParticipantsCmd.AddCommand(listBaselineWorkgroupParticipantsCmd)
}
