package stack

import (
	"log"
	"os"
	"sync"

	"github.com/docker/docker/client"
	"github.com/provideplatform/provide-cli/cmd/common"

	"github.com/spf13/cobra"
)

var logsBaselineStackCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print baseline stack logs",
	Long:  `Print the logs from each container in a local baseline stack instance`,
	Run:   stackLogs,
}

func stackLogs(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepLogs)
}

func stackLogsRun(cmd *cobra.Command, args []string) {
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	common.LogContainers(docker, &wg, name)
	wg.Wait()
}

func init() {
	logsBaselineStackCmd.Flags().StringVar(&name, "name", "baseline-local", "name of the baseline stack instance")
}
