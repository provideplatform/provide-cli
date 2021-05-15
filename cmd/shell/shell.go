package shell

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const shellExitMessage = "Bye!"
const shellTitle = "prvd"
const shellPrefix = ">>> "
const shellOptionInputTextColor = prompt.Green

var ShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interactive shell",
	Long: fmt.Sprintf(`%s

The Provide shell allows you to attach to a specific version of the Provide stack.

Run with the --help flag to see available options`, common.ASCIIBanner),
	Run: shell,
}

func shell(cmd *cobra.Command, args []string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Caught error exception: %v", r)
		}
	}()

	defer fmt.Println(shellExitMessage)

	p := prompt.New(
		func(selected string) {
		},
		completer,
		prompt.OptionTitle(shellTitle),
		prompt.OptionPrefix(shellPrefix),
		prompt.OptionInputTextColor(shellOptionInputTextColor),
	)
	p.Run()
}
