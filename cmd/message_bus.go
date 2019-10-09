package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var subject string

var messageBusCmd = &cobra.Command{
	Use:   "message_bus",
	Short: "Manage message bus applications",
	Long:  `Manage message bus applications consisting of a load balanced distributed filesystem (i.e., an IPFS connector) and an on-chain registry smart contract`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("message bus unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(messageBusCmd)
	messageBusCmd.AddCommand(messageBusInitCmd)
	messageBusCmd.AddCommand(messageBusPublishCmd)
}
