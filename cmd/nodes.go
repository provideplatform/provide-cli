package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var node map[string]interface{}
var nodes []interface{}
var nodeID string
var nodeType string

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manage nodes",
	Long:  `Manage and provision elastic distributed nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nodes init")
	},
}

func init() {
	rootCmd.AddCommand(nodesCmd)
	nodesCmd.AddCommand(nodesInitCmd)
}
