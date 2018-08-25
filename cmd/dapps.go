package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dappsCmd represents the dapps command
var dappsCmd = &cobra.Command{
	Use:   "dapps",
	Short: "Access dapp-specific functionality made available by the provide API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dapps called")
	},
}

func init() {
	rootCmd.AddCommand(dappsCmd)
}
