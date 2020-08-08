package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const asciiBanner = `
██████╗ ██████╗  ██████╗ ██╗   ██╗██╗██████╗ ███████╗
██╔══██╗██╔══██╗██╔═══██╗██║   ██║██║██╔══██╗██╔════╝
██████╔╝██████╔╝██║   ██║██║   ██║██║██║  ██║█████╗  
██╔═══╝ ██╔══██╗██║   ██║╚██╗ ██╔╝██║██║  ██║██╔══╝  
██║     ██║  ██║╚██████╔╝ ╚████╔╝ ██║██████╔╝███████╗
╚═╝     ╚═╝  ╚═╝ ╚═════╝   ╚═══╝  ╚═╝╚═════╝ ╚══════╝`

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "prvd",
	Short: "Provide CLI",
	Long: fmt.Sprintf(`%s

The Provide CLI exposes convenient tools to manage network and application resources.

Run with the --help flag to see available options`, asciiBanner),
}

// Execute the default command path
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}
