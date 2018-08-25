package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dappsCmd represents the dapps command
var dappsCmd = &cobra.Command{
	Use:   "dapps",
	Short: "Access dapp-specific functionality available by the provide microservices.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dapps called")
	},
}

func init() {
	rootCmd.AddCommand(dappsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dappsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dappsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
