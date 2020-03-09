package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"
	"github.com/spf13/cobra"
)

var networkName string
var networksInitCmd = &cobra.Command{
	Use:   "init --name Network1 --application 024ff1ef-7369-4dee-969c-1918c6edb5d4",
	Short: "Initialize a new network",
	Long:  `Initialize a new network with options`,
	Run:   CreateNetwork,
}

// CreateNetwork configures a new peer-to-peer network;
// see https://docs.provide.services/microservices/goldmine/#create-a-network
func CreateNetwork(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{
		"name":           networkName,
		"application_id": applicationID,
	}
	status, resp, err := provide.CreateNetwork(token, params)
	if err != nil {
		log.Printf("Failed to initialize network; %s", err.Error())
		os.Exit(1)
	}
	if status == 201 {
		network = resp.(map[string]interface{})
		networkID = contract["id"].(string)
		result := fmt.Sprintf("%s\t%s\n", contract["id"], contract["name"])
		fmt.Print(result)
	}
}

func init() {
	networksInitCmd.Flags().StringVar(&networkName, "name", "", "name of the network")
	networksInitCmd.MarkFlagRequired("name")

	networksInitCmd.Flags().StringVar(&applicationID, "application", "", "ID of the application")
	networksInitCmd.MarkFlagRequired("application")
}
