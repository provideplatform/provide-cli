package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const shellExitMessage = "Bye!"
const shellTitle = "prvd"
const shellPrefix = "➜  prvd: ✗ "

const shellOptionDefaultInputTextColor = prompt.DefaultColor
const shellOptionDefaultPrefixTextColor = prompt.Green
const shellOptionDescriptionBGColor = prompt.White
const shellOptionDescriptionTextColor = prompt.Black
const shellOptionSelectedDescriptionBGColor = prompt.LightGray
const shellOptionSelectedDescriptionTextColor = prompt.Black
const shellOptionScrollBGColor = prompt.LightGray
const shellOptionScrollColor = prompt.DarkGray
const shellOptionSelectedSuggestionBGColor = prompt.LightGray
const shellOptionSelectedSuggestionTextColor = prompt.Black
const shellOptionSuggestionBGColor = prompt.White
const shellOptionSuggestionTextColor = prompt.Black

const sanitizedPromptInputMatchExit = "exit"
const sanitizedPromptInputMatchQuit = "quit" // FIXME-- combine exit and quit into regex i.e. ^(exit|quit)$
const sanitizedPromptInputMatchRoot = "prvd"

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

	var p *prompt.Prompt
	p = prompt.New(
		func(input string) {
			execInput(cmd, p, input)
		},

		func(d prompt.Document) []prompt.Suggest {
			return promptSuggestionFactory(cmd, d)
		},

		prompt.OptionDescriptionBGColor(shellOptionDescriptionBGColor),
		prompt.OptionDescriptionTextColor(shellOptionDescriptionTextColor),
		prompt.OptionInputTextColor(shellOptionDefaultInputTextColor),
		prompt.OptionLivePrefix(func() (string, bool) {
			return shellPrefix, true
		}),
		prompt.OptionPrefix(shellPrefix),
		prompt.OptionPrefixTextColor(shellOptionDefaultPrefixTextColor),
		prompt.OptionScrollbarBGColor(shellOptionScrollBGColor),
		prompt.OptionScrollbarThumbColor(shellOptionScrollColor),
		prompt.OptionSelectedDescriptionBGColor(shellOptionSelectedDescriptionBGColor),
		prompt.OptionSelectedDescriptionTextColor(shellOptionSelectedDescriptionTextColor),
		prompt.OptionSelectedSuggestionBGColor(shellOptionSelectedSuggestionBGColor),
		prompt.OptionSelectedSuggestionTextColor(shellOptionSelectedSuggestionTextColor),
		// prompt.OptionSetExitCheckerOnInput(),
		prompt.OptionSuggestionBGColor(shellOptionSuggestionBGColor),
		prompt.OptionSuggestionTextColor(shellOptionSuggestionTextColor),
		prompt.OptionTitle(shellTitle),
	)

	p.Run()
}

func execInput(cmd *cobra.Command, p *prompt.Prompt, input string) {
	switch input {
	case sanitizedPromptInputMatchExit:
		os.Exit(0)
	case sanitizedPromptInputMatchQuit:
		os.Exit(0)
	}
}

func promptSuggestionFactory(cmd *cobra.Command, d prompt.Document) []prompt.Suggest {
	input := strings.TrimSpace(d.CurrentLine()) // this is hardly sanitized -- but it's a start
	results := make([]prompt.Suggest, 0)

	switch input {
	case sanitizedPromptInputMatchRoot:
		for _, cmd := range cmd.Parent().Commands() {
			results = append(results, prompt.Suggest{
				Text:        cmd.Use,
				Description: cmd.Short,
			})
		}
	}

	return results
}
