package stack

import (
	"context"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/spf13/cobra"
)

var stopBaselineStackCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the baseline stack",
	Long:  `Stop a local baseline stack instance`,
	Run:   stopProxy,
}

func stopProxy(cmd *cobra.Command, args []string) {
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	purgeContainers(docker)
	purgeNetwork(docker)

	log.Printf("%s local baseline instance stopped", name)
}

func purgeContainers(docker *client.Client) {
	// log.Printf("purging containers for local baseline instance: %s", name)
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

func purgeNetwork(docker *client.Client) {
	networks, _ := docker.NetworkList(context.Background(), types.NetworkListOptions{})
	for _, ntwrk := range networks {
		if ntwrk.Name == name {
			docker.NetworkRemove(context.Background(), ntwrk.ID)
		}
	}
}

func init() {
	stopBaselineStackCmd.Flags().StringVar(&name, "name", "baseline-local", "name of the baseline stack instance")
}
