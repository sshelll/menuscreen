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

type MenuItem struct {
	Content string
	Item    any
}

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

// AppendLines append lines to the end of the menu.
// WARN: do not call AppendItems and AppendLines at the same time.
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

// AppendItems append items to the end of the menu.
// WARN: do not call AppendItems and AppendLines at the same time.
func (menu *MenuScreen) AppendItems(items ...*MenuItem) *MenuScreen {
	menu.items = append(menu.items, items...)
	for _, item := range items {
		menu.lines = append(menu.lines, item.Content)
	}
	return menu
}

func (menu *MenuScreen) ChosenItem() (idx int, item *MenuItem, ok bool) {
	if !menu.confirmed {
		return -1, nil, false
	}
	ok = true
	idx = menu.cursorY
	if menu.mode == modeI {
		idx = -1
		item = &MenuItem{
			Content: string(menu.input),
		}
	}
	if menu.mode == modeN && len(menu.lines) > 0 {
		item = menu.items[idx]
	}
	if menu.mode == modeS && len(menu.matchedLns) > 0 {
		m := menu.matchedLns[idx]
		idx = m.idx
		item = &MenuItem{
			Content: m.content,
			Item:    m.item,
		}
	}
	return
}

func (menu *MenuScreen) refreshScreen() {

	defer menu.screen.Show()

	switch menu.mode {
	case modeN:
		// PERF: only fill screen when the screen is not filled.
		if !menu.hasFilled {
			menu.screen.Clear()
			menu.fillScreen(strSliceToMatchedLines(menu.lines))
			menu.hasFilled = true
		}
		menu.screen.HideCursor()
		menu.resetChosenLine()
	case modeS:
		menu.screen.Clear()
		lns := make([]*matchedLine, 0, 16)
		lns = append(lns, &matchedLine{
			content: "/" + string(menu.query),
		})
		lns = append(lns, menu.matchedLns...)
		menu.fillScreen(lns)
		menu.resetChosenLine()
		cell := cellCnt(menu.query[:menu.inputCursorPos])
		menu.screen.ShowCursor(cell+3, 1)
	case modeI:
		menu.screen.Clear()
		lns := append([]string{":" + string(menu.input)}, menu.lines...)
		menu.fillScreen(strSliceToMatchedLines(lns))
		cell := cellCnt(menu.input[:menu.inputCursorPos])
		menu.screen.ShowCursor(cell+3, 1)
	}

}

// resetChosenLine re-draw the chosen line with cursor.
func (menu *MenuScreen) resetChosenLine() {

	if menu.mode == modeN {
		// reset last chosen line
		menu.setLineWithStyle(menu.lastCursorY+1, "  "+menu.lines[menu.lastCursorY], nil, defaultContentStyle)
		// highlight the chosen line
		menu.setLineWithStyle(menu.cursorY+1, "  "+menu.lines[menu.cursorY], nil, defaultChosenLineStyle)
		// draw the cursor arrow
		menu.setRuneOfLine(0, menu.cursorY+1, '▸', defaultChosenLineStyle)
		return
	}

	if menu.mode == modeS {
		if len(menu.matchedLns) > 0 && menu.cursorY < len(menu.matchedLns) {
			ln := menu.matchedLns[menu.cursorY]
			menu.setLineWithStyle(2+menu.cursorY, "  "+ln.content, ln.pos, defaultChosenLineStyle)
			menu.setRuneOfLine(0, menu.cursorY+2, '▸', defaultChosenLineStyle)
		}
		return
	}

}

func (menu *MenuScreen) fillScreen(lines matchedLines) {

	// title
	menu.setLineWithStyle(0, menu.title, nil, defaultTitleStyle)

	// content
	for i, ln := range lines {
		style := defaultContentStyle
		if i == 0 && menu.mode == modeS {
			style = defaultQueryStyle
		}
		menu.setLineWithStyle(i+1, "  "+ln.content, ln.pos, style)
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
	menu.setLineWithStyle(len(lines)+1, statistic, nil, defaultContentStyle)

}

func (menu *MenuScreen) setLineWithStyle(y int, content string, hlPos []int, style tcell.Style) {
	x := 0
	pset := make(map[int]struct{})
	for _, p := range hlPos {
		pset[p] = struct{}{}
	}
	pos := 0
	for _, c := range content {
		r, w, comb := menu.calRuneWidthAndComb(c)
		targetStyle := style
		if _, ok := pset[pos-2]; ok {
			targetStyle = defaultHighlightStyle
		}
		menu.screen.SetContent(x, y, r, comb, targetStyle)
		x += w
		pos++
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
