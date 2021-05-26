package messages

import (
	"github.com/spf13/cobra"
)

var MessagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Interact with a baseline workflows",
	Long:  `Create, manage and interact with workflows via the baseline protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		generalPrompt(cmd, args, "")
	},
}

func init() {
	MessagesCmd.AddCommand(listBaselineMessagesCmd)
	MessagesCmd.AddCommand(sendBaselineMessageCmd)
	MessagesCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
