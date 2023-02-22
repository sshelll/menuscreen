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
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func (menu *MenuScreen) SetTitle(title string) *MenuScreen {
	menu.title = title
	return menu
}

func (menu *MenuScreen) ClearLines() *MenuScreen {
	menu.lines = nil
	menu.matchedLns = nil
	return menu
}

func (menu *MenuScreen) SetLines(lns ...string) *MenuScreen {
	menu.cursorY = 0
	menu.lines = lns
	return menu
}

func (menu *MenuScreen) SetLine(n int, content string) *MenuScreen {
	menu.cursorY = 0
	if n < 0 {
		panic("line number should greater than 0")
	} else if n < len(menu.lines) {
		menu.lines[n] = content
	} else if n == len(menu.lines) {
		menu.lines = append(menu.lines, content)
	} else if n > len(menu.lines)-1 {
		newLines := make([]string, n+1)
		copy(newLines, menu.lines)
		newLines[n] = content
		menu.lines = newLines
	}
	return menu
}

func (menu *MenuScreen) AppendLines(content ...string) *MenuScreen {
	menu.lines = append(menu.lines, content...)
	return menu
}

func (menu *MenuScreen) ChosenLine() (idx int, ln string, ok bool) {
	if !menu.confirmed {
		return -1, "", false
	}
	ok = true
	idx = menu.cursorY
	if menu.mode == modeI {
		idx = -1
		ln = string(menu.input)
	}
	if menu.mode == modeN && len(menu.lines) > 0 {
		ln = menu.lines[idx]
	}
	if menu.mode == modeS && len(menu.matchedLns) > 0 {
		ln = menu.matchedLns[idx].content
		idx = menu.matchedLns[idx].idx
	}
	return
}

func (menu *MenuScreen) refreshScreen() {

	menu.screen.Clear()
	defer menu.screen.Show()

	if menu.mode == modeN {
		menu.screen.HideCursor()
		menu.fillScreen(menu.lines)
		menu.resetChosenLine()
		return
	}

	if menu.mode == modeS {
		lns := append([]string{"/" + string(menu.query)}, menu.matchedLns.Content()...)
		menu.fillScreen(lns)
		menu.resetChosenLine()
		cell := cellCnt(menu.query[:menu.inputCursorPos])
		menu.screen.ShowCursor(cell+3, 1)
		return
	}

	if menu.mode == modeI {
		lns := append([]string{":" + string(menu.input)}, menu.lines...)
		menu.fillScreen(lns)
		cell := cellCnt(menu.input[:menu.inputCursorPos])
		menu.screen.ShowCursor(cell+3, 1)
		return
	}

}

// resetChosenLine re-draw the chosen line with cursor.
func (menu *MenuScreen) resetChosenLine() {

	if menu.mode == modeN {
		// highlight the chosen line
		menu.setLineWithStyle(menu.cursorY+1, "  "+menu.lines[menu.cursorY], defaultChosenLineStyle)
		// draw the cursor arrow
		menu.setRuneOfLine(0, menu.cursorY+1, '▸', defaultChosenLineStyle)
		return
	}

	if menu.mode == modeS {
		if len(menu.matchedLns) > 0 {
			menu.setLineWithStyle(2+menu.cursorY, "  "+menu.matchedLns[menu.cursorY].content, defaultChosenLineStyle)
			menu.setRuneOfLine(0, menu.cursorY+2, '▸', defaultChosenLineStyle)
		}
		return
	}

}

func (menu *MenuScreen) fillScreen(lines []string) {

	// title
	menu.setLineWithStyle(0, menu.title, defaultTitleStyle)

	// content
	for i, ln := range lines {
		style := defaultContentStyle
		if i == 0 && menu.mode == modeS {
			style = defaultQueryStyle
		}
		menu.setLineWithStyle(i+1, "  "+ln, style)
		// highlight the cursor column
		if menu.mode == modeS && i == 0 {
			continue
		}
		menu.screen.SetContent(0, i+1, ' ', nil, defaultCursorColStyle)
	}

	// statistic
	statistic := fmt.Sprintf("%d/%d", len(menu.lines), len(menu.lines))
	if menu.mode == modeS {
		statistic = fmt.Sprintf("%d/%d", len(menu.matchedLns), len(menu.lines))
	}
	menu.setLineWithStyle(len(lines)+1, statistic, defaultContentStyle)

}

func (menu *MenuScreen) setLineWithStyle(y int, content string, style tcell.Style) {
	x := 0
	for _, c := range content {
		r, w, comb := menu.calRuneWidthAndComb(c)
		menu.screen.SetContent(x, y, r, comb, style)
		x += w
	}
}

func (menu *MenuScreen) setRuneOfLine(x, y int, c rune, style tcell.Style) {
	r, _, comb := menu.calRuneWidthAndComb(c)
	menu.screen.SetContent(x, y, r, comb, style)
}

func (*MenuScreen) calRuneWidthAndComb(c rune) (r rune, width int, comb []rune) {
	width = runewidth.RuneWidth(c)
	if width == 0 {
		comb = []rune{c}
		c = ' '
		width = 1
	}
	r = c
	return
}

func (*MenuScreen) calRuneWidth(s string) (w int) {
	for _, r := range s {
		w += runewidth.RuneWidth(r)
	}
	return
}
