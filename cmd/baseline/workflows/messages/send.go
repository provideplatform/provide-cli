package messages

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var sendBaselineMessageCmd = &cobra.Command{
	Use:   "send",
	Short: "Send baseline message",
	Long:  `Send baseline message in the context of a workflow`,
	Run:   sendMessage,
}

func sendMessage(cmd *cobra.Command, args []string) {
	log.Printf("not implemented")
	os.Exit(1)
}

func init() {

}
