package baseline

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/spf13/cobra"
)

var logsBaselineProxyCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print baseline proxy logs",
	Long:  `Print the logs from each container in a local baseline proxy instance`,
	Run:   logsProxy,
}

func logsProxy(cmd *cobra.Command, args []string) {
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	logContainers(docker, &wg)
	wg.Wait()
}

func logContainers(docker *client.Client, wg *sync.WaitGroup) error {
	for _, container := range listContainers(docker) {
		if wg != nil {
			wg.Add(1)
		}

		containerID := make([]byte, len(container.ID))
		copy(containerID, container.ID)

		go func() {
			logContainer(docker, string(containerID))
			if wg != nil {
				wg.Done()
			}
		}()
	}

	return nil
}

func logContainer(docker *client.Client, containerID string) error {
	out, err := docker.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
	})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return nil
}

func init() {
	logsBaselineProxyCmd.Flags().StringVar(&name, "name", "baseline-proxy", "name of the baseline proxy instance")
}
