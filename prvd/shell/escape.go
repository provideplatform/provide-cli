/*
 * Copyright 2017-2022 Provide Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package shell

func clear() {
	writeRaw([]byte("\033[H\033[2J"), true)
}

func clearScrollback() {
	writeRaw([]byte("\033[3J"), true)
}

func cursorUp() {
	writeRaw([]byte("\033[A"), true)
}

func cursorDown() {
	writeRaw([]byte("\033[B"), true)
}

func cursorRight() {
	writeRaw([]byte("\033[C"), true)
}

func cursorLeft() {
	writeRaw([]byte("\033[D"), true)
}

func cursorNextLine() {
	writeRaw([]byte("\033[E"), true)
}

func cursorPreviousLine() {
	writeRaw([]byte("\033[F"), true)
}

func defaultCursorPosition() {
	if writer != nil {
		writer.CursorGoTo(shellHeaderRows+1, 0)
	}
}

func eraseCursorToEnd() {
	writeRaw([]byte("\033[0J"), true)
}

func eraseCursorToLineEnd() {
	writeRaw([]byte("\033[0K"), true)
}

func eraseCurrentLineToCursor() {
	writeRaw([]byte("\033[1K"), true)
}

func eraseCurrentLine() {
	writeRaw([]byte("\033[2K"), true)
}

func hideCursor() {
	if writer != nil {
		writer.HideCursor()
		cursorHidden = true
	}
}

func saveCursor() {
	if writer != nil {
		writer.SaveCursor()
	}
}

func showCursor() {
	if writer != nil {
		writer.ShowCursor()
		cursorHidden = false
	}
}

func stripEscapeSequences(in string) string {
	if len(in) > 0 {
		// fmt.Printf("%s", in)
		// HACK
		return in[7:]
	}
	return in
}

func toggleAlternateBuffer() {
	if !viewingAlternateBuffer {
		writeRaw([]byte("\033[?1049h"), true)
		viewingAlternateBuffer = true
	} else {
		writeRaw([]byte("\033[?1049l"), true)
		viewingAlternateBuffer = false
	}
}
