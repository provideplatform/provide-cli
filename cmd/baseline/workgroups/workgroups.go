package workgroups

import (
	"fmt"

	"github.com/spf13/cobra"
)

var WorkgroupsCmd = &cobra.Command{
	Use:   "workgroups",
	Short: "Interact with a baseline workgroups",
	Long:  `Create, manage and interact with workgroups via the baseline protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workgroups unimplemented")
	},
}

func init() {
	WorkgroupsCmd.AddCommand(initBaselineWorkgroupCmd)
	WorkgroupsCmd.AddCommand(listBaselineWorkgroupsCmd)
}
