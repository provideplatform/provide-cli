package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const shellExitMessage = "Exiting... have a nice day!"
const shellTitle = "prvd"
const shellPrefix = "➜  prvd:($VERSION)$PATH"
const shellPrefixPrompt = " > "

const shellOptionDefaultFGColor = prompt.DefaultColor
const shellOptionDefaultBGColor = prompt.DefaultColor
const shellOptionDefaultInputTextColor = prompt.DefaultColor
const shellOptionDefaultMaxSuggestions = uint16(8)
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

const sanitizedPromptInputMatchClear = "clear"
const sanitizedPromptInputMatchExit = "exit"
const sanitizedPromptInputMatchQuit = "quit" // FIXME-- combine exit and quit into regex i.e. ^(exit|quit)$
const sanitizedPromptInputMatchRoot = ""

var debug bool

var path string
var version string

var childCommands []*cobra.Command
var childCommandSuggestions []prompt.Suggest

var prmpt *prompt.Prompt
var history *prompt.History
var parser prompt.ConsoleParser
var writer prompt.ConsoleWriter

var shellOutPending bool

var ShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interactive shell",
	Long: fmt.Sprintf(`%s

The Provide shell allows you to attach to a specific version of the Provide stack.

Run with the --help flag to see available options`, common.ASCIIBanner),
	Run: shell,
}

func shell(cmd *cobra.Command, args []string) {
	childCommands = make([]*cobra.Command, 0)
	childCommandSuggestions = make([]prompt.Suggest, 0)
	for _, child := range cmd.Root().Commands() {
		if child != cmd {
			childCommands = append(childCommands, child)
			childCommandSuggestions = append(childCommandSuggestions, prompt.Suggest{
				Text:        child.Use,
				Description: child.Short,
			})
		}
	}

	if version == "" {
		version = "latest" // FIXME- resolve from environment, or disk if this is running in the context of a release `pwd`
	}

	defer fmt.Println(shellExitMessage)
	refresh(cmd, nil)
}

// refresh initializes (or re-initializes) the prompt REPL
func refresh(cmd *cobra.Command, msg []byte) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("WARNING: recovered from panic; %s", r)
			refresh(cmd, []byte(msg))
		}
	}()

	prefix := strings.ReplaceAll(
		strings.ReplaceAll(shellPrefix, "$VERSION", version),
		"$PATH",
		fmt.Sprintf("%s%s", path, shellPrefixPrompt),
	)

	history = prompt.NewHistory() // TODO -- use this?
	parser = prompt.NewStandardInputParser()
	parser.Setup()

	writer = prompt.NewStdoutWriter()
	renderRootBanner()

	if msg != nil && len(msg) > 0 {
		writer.Write(msg)
		writer.WriteStr("\n\n")
	}

	prmpt = prompt.New(
		func(input string) {
			interpret(cmd, input)
		},

		func(d prompt.Document) []prompt.Suggest {
			return promptSuggestionFactory(cmd, d)
		},

		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(buf *prompt.Buffer) {
				writer.WriteStr("Interrupt\n")
			},
		}),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlD,
			Fn: func(buf *prompt.Buffer) {
				os.Exit(0)
			},
		}),
		prompt.OptionBreakLineCallback(func(d *prompt.Document) {
			// noop
		}),
		prompt.OptionDescriptionBGColor(shellOptionDescriptionBGColor),
		prompt.OptionDescriptionTextColor(shellOptionDescriptionTextColor),
		prompt.OptionInputTextColor(shellOptionDefaultInputTextColor),
		prompt.OptionLivePrefix(func() (string, bool) {
			return prefix, true
		}),
		prompt.OptionMaxSuggestion(shellOptionDefaultMaxSuggestions),
		prompt.OptionParser(parser),
		prompt.OptionPrefix(prefix),
		prompt.OptionPrefixTextColor(shellOptionDefaultPrefixTextColor),
		prompt.OptionScrollbarBGColor(shellOptionScrollBGColor),
		prompt.OptionScrollbarThumbColor(shellOptionScrollColor),
		prompt.OptionSelectedDescriptionBGColor(shellOptionSelectedDescriptionBGColor),
		prompt.OptionSelectedDescriptionTextColor(shellOptionSelectedDescriptionTextColor),
		prompt.OptionSelectedSuggestionBGColor(shellOptionSelectedSuggestionBGColor),
		prompt.OptionSelectedSuggestionTextColor(shellOptionSelectedSuggestionTextColor),
		prompt.OptionSetExitCheckerOnInput(func(in string, breakline bool) bool {
			return false
		}),
		prompt.OptionSuggestionBGColor(shellOptionSuggestionBGColor),
		prompt.OptionSuggestionTextColor(shellOptionSuggestionTextColor),
		prompt.OptionTitle(shellTitle),
		prompt.OptionWriter(writer),
	)

	prmpt.Run()
}

func renderRootBanner() {
	if writer != nil {
		writer.WriteRawStr("\033[H\033[2J")
		writer.SetColor(prompt.Cyan, shellOptionDefaultBGColor, true)
		writer.WriteStr(common.ASCIIBanner)
		writer.WriteStr("\n\n")
		writer.SetColor(shellOptionDefaultFGColor, shellOptionDefaultBGColor, true)
	}
}

func shellOut(argv []string) error {
	if shellOutPending {
		return nil
	}

	shellOutPending = true

	buf := &bytes.Buffer{}
	cmd := exec.Command("prvd", argv...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = nil
	cmd.Stdout = buf

	go func() {
		for shellOutPending {
			if buf.Len() > 0 {
				writer.WriteRawStr("\033[2K")
				writer.WriteRawStr(buf.String())
				writer.Flush()
				buf.Truncate(0)
			}

			time.Sleep(time.Millisecond * 50)
		}
	}()

	err := cmd.Run()
	shellOutPending = false
	if err != nil {
		return err
	}

	return nil
}

func interpret(cmd *cobra.Command, input string) {
	if strings.TrimSpace(input) == "" {
		return
	}

	switch input {
	case sanitizedPromptInputMatchClear:
		refresh(cmd, []byte{})
		return
	case sanitizedPromptInputMatchExit:
		os.Exit(0)
		return
	case sanitizedPromptInputMatchQuit:
		os.Exit(0)
		return
	}

	argv := strings.Split(strings.TrimSpace(input), " ")

	_cmd, i := resolveChildCmd(cmd, argv)
	if _cmd != nil {
		if debug {
			writer.WriteStr(fmt.Sprintf("resolved child command for input: %s; argv[%d]: %v; use: %s", input, i, argv, _cmd.Use))
		}

		// // TODO? -- use the following instead of shellOut()
		// _cmd.SetArgs(argv[i:])
		// out := &bytes.Buffer{}
		// _cmd.SetErr(out)
		// _cmd.SetOut(out)
		// _cmd.SetOutput(out)
		// _cmd.SetIn(os.Stdin)
		// err := _cmd.RunE(_cmd, argv[i:])
		// if err != nil {
		// 	writer.WriteStr(fmt.Sprintf("WARNING: failed to execute; %s", err.Error()))
		// }
		// writer.WriteStr(out.String())

		err := shellOut(argv)
		if err != nil {
			panic(err)
		}
	} else {
		writer.WriteStr(fmt.Sprintf("%s: command not found: %s\n", shellTitle, strings.Join(argv, " ")))
	}
}

// resolveChildCmd checks to see if the given input can be resolved to the provide-cli
func resolveChildCmd(cmd *cobra.Command, argv []string) (*cobra.Command, int) {
	for _, child := range childCommands {
		for i := range argv { // FIXME? should this be reversed?
			if strings.TrimSpace(argv[i]) == child.Use {
				return child, i
			}
		}
	}

	return nil, -1
}

func promptSuggestionFactory(cmd *cobra.Command, d prompt.Document) []prompt.Suggest {
	var results []prompt.Suggest = nil

	input := strings.TrimSpace(d.CurrentLine())
	argv := strings.Split(input, " ") // this is hardly sanitized -- but it's a start

	i := len(argv) - 1 // TODO: handle caret position for previous suggest fields...
	// i now contains offset for arg at cursor position

	switch argv[i] {
	case sanitizedPromptInputMatchRoot:
		results = childCommandSuggestions
	}

	return results
}
