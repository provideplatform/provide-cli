package api_tokens

import (
	"os"

	"github.com/spf13/cobra"
)

var APITokensCmd = &cobra.Command{
	Use:   "api_tokens",
	Short: "Manage API tokens",
	Long:  `API tokens can be created on behalf of a developer account, application or application user`,
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
	APITokensCmd.AddCommand(apiTokensListCmd)
	APITokensCmd.AddCommand(apiTokensInitCmd)
	APITokensCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
}
