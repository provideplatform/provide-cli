package shell

import (
	"bytes"
	"sync"
)

const nativeCommandTop = "top"

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
				writer.CursorGoTo(8, 0)
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
