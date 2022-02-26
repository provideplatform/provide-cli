package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/provideplatform/provide-cli/prvd/common"
)

var rootCmd = &cobra.Command{
	Use:   "prvdnetwork",
	Short: "provide.network cli",
	Long: fmt.Sprintf(`%s

The provide.network CLI exposes commands to run and manage full and validator nodes on permissionless provide.network.

Run with the --help flag to see available options`, common.ASCIIBanner),
}

// Execute the default command path
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(common.InitConfig)

	rootCmd.PersistentFlags().BoolVarP(&common.Verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&common.CfgFile, "config", "c", "", "config file (default is $HOME/.provide-cli.yaml)")

	common.CacheCommands(rootCmd)
}
