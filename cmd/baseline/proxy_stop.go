package baseline

import (
	"context"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/spf13/cobra"
)

var stopBaselineProxyCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the baseline proxy",
	Long:  `Stop a local baseline proxy instance`,
	Run:   stopProxy,
}

func stopProxy(cmd *cobra.Command, args []string) {
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	purgeContainers(docker)

	log.Printf("%s proxy instance stopped", name)
}

func purgeContainers(docker *client.Client) {
	for _, container := range listContainers(docker) {
		err := docker.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})

		if err != nil {
			log.Printf("WARNING: failed to remove container: %s; %s", container.Names[0], err.Error())
		}
	}
}

func init() {
	stopBaselineProxyCmd.Flags().StringVar(&name, "name", "baseline-proxy", "name of the baseline proxy instance")
}
