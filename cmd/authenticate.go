package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// authenticateCmd represents the authenticate command
var authenticateCmd = &cobra.Command{
	Use:   "prvd authenticate",
	Short: "Authenticate using user credentials for provide.services.",
	Long:  `Authenticate using user credentials for provide.services`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("authenticate called")
	},
}

func init() {
	rootCmd.AddCommand(authenticateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authenticateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authenticateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
