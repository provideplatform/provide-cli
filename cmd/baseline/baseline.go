package baseline

import (
	"fmt"

	"github.com/spf13/cobra"
)

var BaselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Interact with the baseline protocol",
	Long: `Create an manage vaults and their associated keys and secrets.

Supports encrypt/decrypt and sign/verify operations for select key specs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("baseline unimplemented")
	},
}

func init() {
	BaselineCmd.AddCommand(logsBaselineProxyCmd)
	BaselineCmd.AddCommand(runBaselineProxyCmd)
	BaselineCmd.AddCommand(stopBaselineProxyCmd)
}
