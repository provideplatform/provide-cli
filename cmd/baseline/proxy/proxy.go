package proxy

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ProxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Interact with a local baseline proxy",
	Long:  `Create, manage and interact with local baseline proxy instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy unimplemented")
	},
}

func init() {
	ProxyCmd.AddCommand(logsBaselineProxyCmd)
	ProxyCmd.AddCommand(runBaselineProxyCmd)
	ProxyCmd.AddCommand(stopBaselineProxyCmd)
}
