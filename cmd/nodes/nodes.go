package nodes

import (
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var node map[string]interface{}
var nodes []interface{}
var nodeType string

var NodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manage nodes",
	Long:  `Manage and provision elastic distributed nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")

		defer func() {
			if r := recover(); r != nil {
				os.Exit(1)
			}
		}()
	},
}

func init() {
	NodesCmd.AddCommand(nodesInitCmd)
	NodesCmd.AddCommand(nodesLogsCmd)
	NodesCmd.AddCommand(nodesDeleteCmd)
	NodesCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")

}
