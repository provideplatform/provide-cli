package networks

import (
	"fmt"

	"github.com/spf13/cobra"
)

var network map[string]interface{}
var networks []interface{}

var NetworksCmd = &cobra.Command{
	Use:   "networks",
	Short: "Manage networks",
	Long:  `Manage and provision elastic distributed networks and other infrastructure`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("networks unimplemented")
	},
}

func init() {
	NetworksCmd.AddCommand(networksInitCmd)
	NetworksCmd.AddCommand(networksListCmd)
	NetworksCmd.AddCommand(networksDisableCmd)
}
