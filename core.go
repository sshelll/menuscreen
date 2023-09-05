/*
 * Copyright (c) 2021. shaojiale
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use file except in compliance with the License.
 * You may obtain a copy of the license at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package menuscreen

import (
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// This file includes all reserved key bind mapping.
// '↑ ↓' controls the cursor;
// 'esc' means exit the current mode;
// 'enter' means a line has been chosen;
// 'slash' means enter the query mode;
// 'runes' means input in the query mode;
// 'backspace' means rollback the last char from input;

func (menu *MenuScreen) keyUP(*tcell.EventKey) {
	if menu.cursorY > 0 {
		menu.cursorY--
	} else if menu.mode == modeN {
		menu.cursorY = len(menu.lines) - 1
	} else if menu.mode == modeS {
		menu.cursorY = len(menu.matchedLns) - 1
	}
}

func (menu *MenuScreen) keyDOWN(*tcell.EventKey) {
	if menu.checkCursor() {
		menu.cursorY++
	} else {
		menu.cursorY = 0
	}
}

func (menu *MenuScreen) keyRIGHT(*tcell.EventKey) {
	if menu.mode == modeI {
		menu.inputCursorPos = min(menu.inputCursorPos+1, len(menu.input))
	}
	if menu.mode == modeS {
		menu.inputCursorPos = min(menu.inputCursorPos+1, len(menu.query))
	}
}

func (menu *MenuScreen) keyLEFT(*tcell.EventKey) {
	if menu.mode == modeI || menu.mode == modeS {
		menu.inputCursorPos = max(menu.inputCursorPos-1, 0)
	}
}

func (menu *MenuScreen) keyESC(*tcell.EventKey) {
	switch menu.mode {
	case modeS, modeI:
		menu.mode = modeN
		menu.inputCursorPos = 0
	default:
		menu.shutdown()
	}
}

func (menu *MenuScreen) keyENTER(*tcell.EventKey) {
	menu.confirmed = true
	menu.inputCursorPos = 0
	menu.shutdown()
}

func (menu *MenuScreen) keyBS(*tcell.EventKey) {

	if menu.mode == modeI {
		if len(menu.input) == 0 {
			return
		}
		newRunes := []rune(menu.input)
		delPos := max(0, menu.inputCursorPos-1)
		delPos = min(delPos, len(newRunes)-1)
		newRunes = append(newRunes[:delPos], newRunes[delPos+1:]...)
		menu.input = newRunes
		menu.inputCursorPos = max(menu.inputCursorPos-1, 0)
		return
	}

	if menu.mode == modeS {
		if len(menu.query) == 0 {
			return
		}
		newRunes := []rune(menu.query)
		delPos := max(0, menu.inputCursorPos-1)
		delPos = min(delPos, len(newRunes)-1)
		newRunes = append(newRunes[:delPos], newRunes[delPos+1:]...)
		menu.query = newRunes
		menu.inputCursorPos = max(menu.inputCursorPos-1, 0)
		menu.cursorY = 0
		menu.calMatchedLines()
	}

}

// keyRUNE controls the input of user.
func (menu *MenuScreen) keyRUNE(ev *tcell.EventKey) {

	runeName := menu.getRuneName(ev.Name())

	if menu.mode == modeS {
		newRunes := append(cloneRuneSlice(menu.query)[:menu.inputCursorPos], []rune(runeName)...)
		newRunes = append(newRunes, []rune(menu.query)[menu.inputCursorPos:]...)
		menu.query = newRunes
		menu.calMatchedLines()
		menu.inputCursorPos = min(menu.inputCursorPos+1, len(menu.query))
		menu.cursorY = 0
		return
	}

	if menu.mode == modeI {
		newRunes := append(cloneRuneSlice(menu.input)[:menu.inputCursorPos], []rune(runeName)...)
		newRunes = append(newRunes, []rune(menu.input)[menu.inputCursorPos:]...)
		menu.input = newRunes
		menu.inputCursorPos = min(menu.inputCursorPos+1, len(menu.input))
		return
	}

	if menu.mode == modeN {
		if runeName == slash {
			menu.keySLASH()
		}
		if runeName == colon {
			menu.keyCOLON()
		}
		return
	}

}

// keySLASH key slash make MenuScreen enter search mode.
func (menu *MenuScreen) keySLASH() {
	menu.mode = modeS
	menu.query = nil
	menu.cursorY = 0
	menu.matchedLns = make([]*matchedLine, 0, len(menu.lines))
	for i, ln := range menu.lines {
		menu.matchedLns = append(menu.matchedLns, &matchedLine{
			idx:     i,
			content: ln,
		})
	}
}

// keyCOLON key colon make MenuScreen enter insert mode.
func (menu *MenuScreen) keyCOLON() {
	menu.mode = modeI
	menu.input = nil
	menu.cursorY = 0
}

func (menu *MenuScreen) getRuneName(k string) string {
	if !strings.HasPrefix(k, "Rune") {
		return k
	}
	return strings.ReplaceAll(strings.ReplaceAll(k, "]", ""), "Rune[", "")
}

func (menu *MenuScreen) checkCursor() bool {
	if menu.mode == modeN {
		return menu.cursorY < len(menu.lines)-1
	}
	return menu.cursorY < len(menu.matchedLns)-1
}

func (menu *MenuScreen) calMatchedLines() {
	matched := make([]*matchedLine, 0, len(menu.lines))
	if len(menu.query) == 0 {
		for _, c := range menu.lines {
			matched = append(matched, &matchedLine{content: c})
		}
		return
	}
	results := menu.fuzzyFinder.Clear().AppendTargets(menu.lines...).MergeMatch(string(menu.query))
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score() < results[j].Score()
	})
	for i, ln := range results {
		mln := &matchedLine{
			idx:     i,
			content: ln.Content(),
			pos:     ln.Pos(),
		}
		matched = append(matched, mln)
	}
	menu.matchedLns = matched
}

func (menu *MenuScreen) initKeyBinder() {

	menu.keyBinder = new(keyBinder)

	menu.keyBinder.bind(menu.keyUP, tcell.KeyUp)
	menu.keyBinder.bind(menu.keyDOWN, tcell.KeyDown)
	menu.keyBinder.bind(menu.keyENTER, tcell.KeyEnter)
	menu.keyBinder.bind(menu.keyRUNE, tcell.KeyRune)
	menu.keyBinder.bind(menu.keyBS, tcell.KeyBackspace, tcell.KeyDEL, tcell.KeyDelete)
	menu.keyBinder.bind(menu.keyESC, tcell.KeyEsc)
	menu.keyBinder.bind(menu.keyLEFT, tcell.KeyLeft)
	menu.keyBinder.bind(menu.keyRIGHT, tcell.KeyRight)

}
