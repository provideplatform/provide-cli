package users

import (
	"fmt"

	"github.com/provideplatform/provide-cli/cmd/common"

	"github.com/spf13/cobra"
)

// ShowIDCmd represents the id command
var showIDCmd = &cobra.Command{
	Use:   "id",
	Short: "Prints out the ID of the currently authenticated user",
	Long:  "Prints out the ID of the currently authenticated user",
	Run:   showid,
}

func showid(cmd *cobra.Command, args []string) {
	id := common.DecodeUserID()
	fmt.Println(id)
}
