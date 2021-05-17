package connectors

import (
	"os"

	"github.com/spf13/cobra"
)

const connectorTypeIPFS = "ipfs"

var connector map[string]interface{}
var connectors []interface{}

var ConnectorsCmd = &cobra.Command{
	Use:   "connectors",
	Short: "Manage arbitrary external infrastructure",
	Long: `Connectors are adapters that connect external arbitrary infrastructure with Provide.

This API allows you to provision load balanced, cloud-agnostic infrastructure for your distributed system.`,
	Run: func(cmd *cobra.Command, args []string) {
		generalPrompt(cmd, args, "")

		defer func() {
			if r := recover(); r != nil {
				os.Exit(1)
			}
		}()
	},
}

func init() {
	ConnectorsCmd.AddCommand(connectorsListCmd)
	ConnectorsCmd.AddCommand(connectorsInitCmd)
	ConnectorsCmd.AddCommand(connectorsDetailsCmd)
	ConnectorsCmd.AddCommand(connectorsDeleteCmd)
}
