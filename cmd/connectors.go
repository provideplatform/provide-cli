package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var connectorsCmd = &cobra.Command{
	Use:   "connectors",
	Short: "Manage application connectors",
	Long:  `Connectors are a powerful way to provision load balanced infrastructure for your application, such as a public or private IPFS cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("connectors unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(connectorsCmd)
	connectorsCmd.AddCommand(connectorsListCmd)
}
