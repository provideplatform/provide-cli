package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var connector map[string]interface{}
var connectorID string

var connectorsCmd = &cobra.Command{
	Use:   "connectors",
	Short: "Manage application connectors",
	Long:  `Provision load balanced infrastructure for your application, such as a public or private IPFS cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("connectors unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(connectorsCmd)
	connectorsCmd.AddCommand(connectorsListCmd)
	connectorsCmd.AddCommand(connectorsInitCmd)
	connectorsCmd.AddCommand(connectorsDetailsCmd)
	connectorsCmd.AddCommand(connectorsDeleteCmd)
}
