package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var apiTokensCmd = &cobra.Command{
	Use:   "api_tokens",
	Short: "Manage user and dapp API tokens",
	Long:  `API tokens can be created on behalf of developer accounts or dapps`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("api_tokens unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(apiTokensCmd)
	apiTokensCmd.AddCommand(apiTokensListCmd)
}
