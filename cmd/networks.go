package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var networksCmd = &cobra.Command{
	Use:   "networks",
	Short: "Access network- and devops-specific functionality made available by the provide API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("networks unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(networksCmd)
	networksCmd.AddCommand(networksListCmd)
}
