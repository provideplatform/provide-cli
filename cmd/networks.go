package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var networkID string

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
	networksCmd.AddCommand(networksListCmd)
}
