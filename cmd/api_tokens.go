package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var apiTokensCmd = &cobra.Command{
	Use:   "api_tokens",
	Short: "Manage API tokens",
	Long:  `API tokens can be created on behalf of a developer account, application or application user`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("api_tokens unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(apiTokensCmd)
	apiTokensCmd.AddCommand(apiTokensListCmd)
	apiTokensCmd.AddCommand(apiTokensInitCmd)
}
