package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dappsCmd = &cobra.Command{
	Use:   "dapps",
	Short: "Access dapp-specific functionality made available by the provide API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dapps unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(dappsCmd)
	dappsCmd.AddCommand(dappsListCmd)
	dappsCmd.AddCommand(dappsInitCmd)
}
