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
	"github.com/gdamore/tcell/v2"
	"strings"
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
	}
}

func (menu *MenuScreen) keyDOWN(*tcell.EventKey) {
	if menu.checkCursor() {
		menu.cursorY++
	}
}

func (menu *MenuScreen) keyESC(*tcell.EventKey) {
	if menu.mode == modeS {
		menu.mode = modeN
	} else if menu.mode == modeN {
		menu.Shutdown()
	}
}

func (menu *MenuScreen) keyENTER(*tcell.EventKey) {
	menu.confirmed = true
	menu.Shutdown()
}

func (menu *MenuScreen) keyBS(*tcell.EventKey) {

	if menu.mode != modeS {
		return
	}

	if len(menu.input) == 0 {
		return
	}

	rs := []rune(menu.input)
	rs = rs[:len(rs)-1]
	menu.input = string(rs)

	menu.calMatchedLines()

}

// keyRUNE controls the input of user.
func (menu *MenuScreen) keyRUNE(ev *tcell.EventKey) {

	runeName := menu.getRuneName(ev.Name())

	if menu.mode == modeS {
		menu.input += runeName
		menu.calMatchedLines()
		return
	}

	if menu.mode == modeN {
		if runeName == slash {
			menu.keySLASH()
		}
		return
	}

}

// keySLASH key slash make MenuScreen enter search mode.
func (menu *MenuScreen) keySLASH() {
	menu.mode = modeS
	menu.input = ""
	menu.cursorY = 0
	menu.matchedLns = make([]*matchedLine, 0, len(menu.lines))
	for i, ln := range menu.lines {
		menu.matchedLns = append(menu.matchedLns, &matchedLine{i, ln})
	}
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
	for i, ln := range menu.lines {
		if strings.Contains(ln, menu.input) {
			matched = append(matched, &matchedLine{i, ln})
		}
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

}
