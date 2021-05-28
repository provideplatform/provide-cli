package shell

import (
	"bytes"
	"os"
	"sync"

	"github.com/manifoldco/promptui"
)

const nativeCommandTop = "top"

type NoopCloser struct {
	buf *bytes.Buffer
}

func (nc *NoopCloser) Write(buf []byte) (int, error) {
	// fmt.Printf("write %d-byte buffer...", len(buf))
	// writer.WriteRaw([]byte("\033[0J"))
	// raw := stripEscapeSequences(string(buf))
	// writer.WriteRaw([]byte(raw))
	// // writer.UnSaveCursor()
	// writer.Flush()

	return nc.buf.Write(buf)
}

func (nc *NoopCloser) Close() error {
	return nil
}

var nativeCommands = map[string]interface{}{
	nativeCommandTop: func(argv []string) (*REPL, error) {

		buf := &bytes.Buffer{}
		i := 0
		var pending bool

		saveCursor()
		hideCursor()
		eraseCursorToEnd()

		repl, _ = NewREPL(func(wg *sync.WaitGroup) error {
			if !pending {
				pending = true
				toggleAlternateBuffer()
				go shellOut("docker", []string{"stats"}, buf)
			}

			if buf.Len() > 0 {
				mutex.Lock()
				defer mutex.Unlock()

				writer.SaveCursor()
				writer.CursorGoTo(shellHeaderRows+1, 0)
				writer.WriteRaw([]byte("\033[0J"))
				raw := stripEscapeSequences(string(buf.Bytes()[i:]))
				writer.WriteRaw([]byte(raw))
				writer.UnSaveCursor()
				writer.Flush()

				i += buf.Len() - i
			}

			return nil
		})

		go repl.run()
		return repl, nil
	},
}

// MarshalPromptIO marshals IO from promptui text or select prompt
func MarshalPromptIO(sel *promptui.Select) {
	buf := &bytes.Buffer{}
	sel.Stdin = os.Stdin
	// sel.Stdout = &NoopCloser{
	// 	buf: buf,
	// }

	i := 0
	repl, _ := NewREPL(func(wg *sync.WaitGroup) error {
		if buf.Len() > 0 {
			mutex.Lock()
			defer mutex.Unlock()

			// writer.SaveCursor()
			// writer.HideCursor()
			eraseCurrentLine()
			// writer.CursorGoTo(shellHeaderRows+1, 0)
			writer.WriteRaw([]byte("\033[0J"))
			raw := stripEscapeSequences(string(buf.Bytes()[i:]))
			writer.WriteRaw([]byte(raw))
			// writer.UnSaveCursor()
			// writer.ShowCursor()
			writer.Flush()

			i += buf.Len() - i
		}

		return nil
	})

	// TODO-- don't leak this...
	go repl.run()
}

func resolveNativeCommand(argv []string) func([]string) (*REPL, error) {
	if len(argv) > 0 {
		if fn, fnOk := nativeCommands[argv[0]].(func([]string) (*REPL, error)); fnOk {
			return fn
		}
	}
	return nil
}

func supportedNativeCommand(argv []string) bool {
	if len(argv) > 0 {
		return resolveNativeCommand(argv) != nil
	}
	return false
}
