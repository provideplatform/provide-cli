package node

import (
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/provideplatform/provide-cli/cmd/common"

	"github.com/spf13/cobra"
)

var stopBaseledgerNodeCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a baseledger node",
	Long:  `Stop a local baseledger node`,
	Run:   stopBaseledgerNode,
}

func stopBaseledgerNode(cmd *cobra.Command, args []string) {
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	common.PurgeContainers(docker, name)
	common.PurgeNetwork(docker, name)

	log.Printf("%s local baseledger node stopped", name)
}

func init() {
	stopBaseledgerNodeCmd.Flags().StringVar(&name, "name", "baseledger-local", "name of the baseledger node instance")
}
