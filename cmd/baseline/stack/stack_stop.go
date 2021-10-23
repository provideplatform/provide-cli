package stack

import (
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/provideplatform/provide-cli/cmd/common"

	"github.com/spf13/cobra"
)

var stopBaselineStackCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the baseline stack",
	Long:  `Stop a local baseline stack instance`,
	Run:   stopStack,
}

func stopStack(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepStop)
}

func runStackStop(cmd *cobra.Command, args []string) {
	docker, err := client.NewClientWithOpts()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	common.PurgeContainers(docker, name)
	common.PurgeNetwork(docker, name)

	log.Printf("%s local baseline instance stopped", name)
}

func init() {
	stopBaselineStackCmd.Flags().StringVar(&name, "name", "baseline-local", "name of the baseline stack instance")
}
