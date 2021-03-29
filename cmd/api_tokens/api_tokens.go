package api_tokens

import (
	"fmt"

	"github.com/spf13/cobra"
)

var APITokensCmd = &cobra.Command{
	Use:   "api_tokens",
	Short: "Manage API tokens",
	Long:  `API tokens can be created on behalf of a developer account, application or application user`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("api_tokens unimplemented")
	},
}

func init() {
	APITokensCmd.AddCommand(apiTokensListCmd)
	APITokensCmd.AddCommand(apiTokensInitCmd)
}
