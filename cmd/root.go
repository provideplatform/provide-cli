package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/provideservices/provide-cli/cmd/root"
)

var rootCmd = &cobra.Command{
	Use:   "prvd",
	Short: "Provide CLI",
	Long: fmt.Sprintf(`%s

The Provide CLI exposes low-code tools to manage network, application and organization resources.

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
	root.Root(rootCmd, true)
}
