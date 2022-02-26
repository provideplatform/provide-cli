package messages

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var listBaselineMessagesCmd = &cobra.Command{
	Use:   "list",
	Short: "List baseline messages",
	Long:  `List baseline messages in the context of a workflow`,
	Run:   listMessages,
}

func listMessages(cmd *cobra.Command, args []string) {
	log.Printf("not implemented")
	os.Exit(1)
}

func init() {

}
