package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const shellHeaderStartRow = 1
const shellHeaderRows = 7

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

var cursorHidden bool
var viewingAlternateBuffer bool

var prmpt *prompt.Prompt
var history *prompt.History
var parser prompt.ConsoleParser
var writer prompt.ConsoleWriter

var mutex *sync.Mutex
var repl *REPL
var repls []*REPL
var wg *sync.WaitGroup

var ShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interactive shell",
	Long: fmt.Sprintf(`%s

The Provide shell allows you to attach to a specific version of the Provide stack.

Run with the --help flag to see available options`, common.ASCIIBanner),
	Run: shell,
}

func shell(cmd *cobra.Command, args []string) {
	mutex = &sync.Mutex{}
	wg = &sync.WaitGroup{}

	repls = make([]*REPL, 0)

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
		if common.IsReleaseContext() {
			version = common.Manifest.Version
		} else {
			version = "latest"
		}
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

	clear()
	clearScrollback()
	renderRootBanner()
	defaultCursorPosition()

	if msg != nil && len(msg) > 0 {
		write(msg, false)
		writeRaw([]byte("\n"), true)
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
				if viewingAlternateBuffer {
					for _, repl := range repls { // FIXME-- should this just be a single, top-level repl?
						repl.shutdown()
					}

					repls = make([]*REPL, 0)
					toggleAlternateBuffer()
					showCursor()
				} else {
					writer.WriteRaw([]byte("Interrupt\n"))
				}
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
			if cursorHidden {
				return "", true
			}
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

	installREPL()

	prmpt.Run()
}

func installREPL() {
	repl, _ = NewREPL(func(_wg *sync.WaitGroup) error {
		renderRootBanner()
		return nil
	})
	go repl.run()

	// buf := &bytes.Buffer{}
	// n := 0
	// go shellOut("docker", []string{"stats"}, buf)
	// repl, _ = NewREPL(func(_wg *sync.WaitGroup) error {
	// 	if buf.Len() > 0 {
	// 		mutex.Lock()
	// 		defer mutex.Unlock()

	// 		writer.SaveCursor()
	// 		writer.HideCursor()

	// 		str := string(buf.Bytes()[n:])
	// 		i := 0

	// 		for _, line := range strings.Split(str, "\n") {
	// 			writer.CursorGoTo(1+i, 64)
	// 			writer.WriteRaw([]byte("\033[0K")) // erase to end of line
	// 			raw := stripEscapeSequences(line)
	// 			writer.WriteRawStr(raw)

	// 			i++
	// 		}

	// 		writer.Flush()
	// 		writer.UnSaveCursor()

	// 		n += buf.Len() - n
	// 		buf.Reset()
	// 	}

	// 	return nil
	// })
	// go repl.run()
}

func renderRootBanner() {
	if writer != nil {
		mutex.Lock()
		defer mutex.Unlock()

		shouldShowCursor := !cursorHidden

		writer.SaveCursor()
		hideCursor()

		// writer.CursorGoTo(shellHeaderStartRow, 0)

		i := 0
		for i < shellHeaderRows {
			writer.CursorGoTo(i, 0)
			writer.WriteRaw([]byte("\033[2K\n")) // delete current line
			i++
		}

		writer.CursorGoTo(shellHeaderStartRow, 0)
		writer.SetColor(prompt.Cyan, shellOptionDefaultBGColor, true)
		writer.WriteStr(common.ASCIIBanner)

		// if common.IsReleaseContext() {
		// 	writer.CursorGoTo(shellHeaderRows-1, int(parser.GetWinSize().Col-uint16(len(common.Manifest.Version))))
		// 	writer.WriteStr(common.Manifest.Version)
		// }

		// render single blank link
		writer.CursorGoTo(shellHeaderRows, 0)
		writer.WriteRaw([]byte("\033[2K\n")) // delete current line
		writer.SetColor(shellOptionDefaultFGColor, shellOptionDefaultBGColor, true)
		writer.Flush()
		writer.UnSaveCursor()

		if shouldShowCursor {
			showCursor()
		}
	}
}

func shellOut(bin string, argv []string, buf *bytes.Buffer) error {
	shellOutPending := true

	write := buf == nil
	if write {
		buf = &bytes.Buffer{}
	}

	cmd := exec.Command(bin, argv...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = nil
	cmd.Stdout = buf

	if write {
		go func() {
			for shellOutPending {
				if buf.Len() > 0 {
					eraseCurrentLine()
					writeRaw(buf.Bytes(), true)
					buf.Reset()
				}

				time.Sleep(time.Millisecond * 50)
			}
		}()
	}

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
		defaultCursorPosition()
		eraseCursorToEnd()
		return
	case sanitizedPromptInputMatchExit:
		showCursor()
		os.Exit(0)
		return
	case sanitizedPromptInputMatchQuit:
		showCursor()
		os.Exit(0)
		return
	}

	argv := strings.Split(strings.TrimSpace(input), " ")

	_cmd, i := resolveChildCmd(cmd, argv)
	if _cmd != nil {
		if debug {
			writer.WriteRaw([]byte(fmt.Sprintf("resolved child command for input: %s; argv[%d]: %v; use: %s", input, i, argv, _cmd.Use)))
		}

		// // TODO? -- use the following instead of shellOut()
		// _cmd.SetArgs(argv[i:])
		// _cmd.SetErr(out)
		// _cmd.SetOut(out)
		// _cmd.SetOutput(out)
		// _cmd.SetIn(os.Stdin)
		// err := _cmd.RunE(_cmd, argv[i:])
		// if err != nil {
		// 	writer.WriteStr(fmt.Sprintf("WARNING: failed to execute; %s", err.Error()))
		// }
		// writer.WriteStr(out.String())

		out := &bytes.Buffer{}
		repl, _ := NewREPLWithCmd(*exec.Command("prvd", argv...), out)
		repl.run()
		repls = append(repls, repl)
		writer.WriteRaw(out.Bytes())
	} else if supportedNativeCommand(argv) {
		mutex.Lock()
		writer.SaveCursor()
		writer.HideCursor()
		writer.Flush()
		mutex.Unlock()

		repl, err := resolveNativeCommand(argv)(argv)
		repls = append(repls, repl)
		if err != nil {
			writer.WriteRaw([]byte(fmt.Sprintf("%s: native command returned err: %s;%s\n", shellTitle, strings.Join(argv, " "), err.Error())))
		}
	} else {
		if len(argv) > 0 {
			writer.WriteRaw([]byte(fmt.Sprintf("%s: command not found: %s\n", shellTitle, strings.Join(argv, " "))))
		}
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
	if cursorHidden {
		return nil
	}

	// var results []prompt.Suggest = nil
	input := strings.TrimSpace(d.CurrentLine())
	// argv := strings.Split(input, " ") // this is hardly sanitized -- but it's a start

	// i := len(argv) - 1 // TODO: handle caret position for previous suggest fields...
	// i now contains offset for arg at cursor position

	// switch argv[i] {
	// case sanitizedPromptInputMatchRoot:
	// 	results = childCommandSuggestions
	// }

	results := make([]prompt.Suggest, 0)
	for _, result := range childCommandSuggestions {
		if len(input) > 0 {
			if strings.HasPrefix(result.Text, input) {
				results = append(results, result)
			}
		} else {
			// Uncomment the following line to enable suggestions by default
			// results = append(results, result)
		}
	}

	return results
}

func write(buf []byte, flush bool) error {
	if writer != nil {
		mutex.Lock()
		defer mutex.Lock()

		writer.Write(buf)
		if flush {
			writer.Flush()
		}
	}

	return nil
}

func writeRaw(buf []byte, flush bool) error {
	if writer != nil {
		mutex.Lock()
		defer mutex.Unlock()

		writer.WriteRaw(buf)
		if flush {
			writer.Flush()
		}
	}

	return nil
}
