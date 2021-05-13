package nodes

import (
	"fmt"
	"os"

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
		fmt.Println("Nodes command run")
		generalPrompt(cmd, args, "")

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Prompt Exit\n")
				os.Exit(1)
			}
		}()
	},
}

func init() {
	NodesCmd.AddCommand(nodesInitCmd)
	NodesCmd.AddCommand(nodesLogsCmd)
	NodesCmd.AddCommand(nodesDeleteCmd)
}
