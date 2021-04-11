package baseline

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/provideservices/provide-cli/cmd/baseline/proxy"
	"github.com/provideservices/provide-cli/cmd/baseline/workgroups"
)

var BaselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Interact with the baseline protocol",
	Long:  `Interact with the baseline protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("baseline unimplemented")
	},
}

func init() {
	BaselineCmd.AddCommand(proxy.ProxyCmd)
	BaselineCmd.AddCommand(workgroups.WorkgroupsCmd)
}
