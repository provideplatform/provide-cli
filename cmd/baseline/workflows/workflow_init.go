package workflows

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var name string
var workgroupID string
var Optional bool

var initBaselineWorkflowCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize baseline workflow",
	Long:  `Initialize and configure a new baseline workflow`,
	Run:   initWorkflow,
}

func initWorkflow(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func initWorkflowRun(cmd *cobra.Command, args []string) {
	log.Printf("not implemented")
	os.Exit(1)
}

func init() {
	initBaselineWorkflowCmd.Flags().StringVar(&name, "name", "", "name of the baseline workflow")
	// initBaselineWorkflowCmd.MarkFlagRequired("name")

	initBaselineWorkflowCmd.Flags().StringVar(&workgroupID, "workgroup", "", "workgroup identifier")
	// initBaselineWorkflowCmd.MarkFlagRequired("workgroup")
	initBaselineWorkflowCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
