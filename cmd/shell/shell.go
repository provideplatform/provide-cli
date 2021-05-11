package shell

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var ShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Provide CLI run in a shell environment",
	Long: fmt.Sprintf(`%s

The Provide CLI exposes low-code tools to manage network, application and organization resources.

Run with the --help flag to see available options`, common.ASCIIBanner),
	Run: shell,
}

func shell(cmd *cobra.Command, args []string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Caught error exception: %v", r)
		}
	}()

	defer fmt.Println("Bye!")
	p := prompt.New(
		func(selected string) {
		},
		completer,
		prompt.OptionTitle("prvd-prompt: interactive provide client"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Yellow),
	)
	p.Run()
}
