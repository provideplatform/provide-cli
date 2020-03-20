package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var network map[string]interface{}
var networks []interface{}
var networkID string
var networkType string

var networksCmd = &cobra.Command{
	Use:   "networks",
	Short: "Manage networks",
	Long:  `Manage and provision elastic distributed networks and other infrastructure`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("networks unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(networksCmd)
	networksCmd.AddCommand(networksInitCmd)
	networksCmd.AddCommand(networksListCmd)
	networksCmd.AddCommand(networksDisableCmd)
}
