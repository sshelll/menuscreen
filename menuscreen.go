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
	"runtime/debug"

	"github.com/gdamore/tcell/v2"
	"github.com/sshelll/fzflib"
)

// MenuScreen is a visible selector for input.
//
// Terminal will show the content below in the normal mode:
//
//	-----------------
//
// | ${Title}:       |
// | ▸ ${1st line}   |
// |   ${2nd line}   |
// |   ${3rd line}   |
//
//	-----------------
//
// Once slash was pressed and MenuScreen entered query mode,
// the screen will be like:
//
//	-----------------------
//
// | ${Title}:             |
// | /query                |
// | ▸ ${1st matched line} |
// |   ${2nd matched line} |
// |   ${3rd matched line} |
//
//	-----------------------
type MenuScreen struct {
	screen         tcell.Screen
	keyBinder      *keyBinder
	shutdownCtrl   chan struct{}
	mode           screenMode
	cursorY        int
	query          []rune
	input          []rune
	inputCursorPos int
	title          string
	lines          []string
	items          []*MenuItem
	matchedLns     matchedLines
	confirmed      bool
	finished       bool
	fuzzyFinder    *fzflib.Fzf
}

func NewMenuScreen() (menuScreen *MenuScreen, err error) {

	screen, err := tcell.NewScreen()
	if err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			screen.Fini()
		}
	}()

	if err = screen.Init(); err != nil {
		return
	}

	screen.SetStyle(defaultContentStyle)
	screen.EnablePaste()
	screen.DisableMouse()
	screen.Clear()

	menu := &MenuScreen{
		screen:      screen,
		lines:       make([]string, 0, 16),
		matchedLns:  make([]*matchedLine, 0, 16),
		mode:        modeN,
		cursorY:     0,
		query:       nil,
		title:       "Menu",
		fuzzyFinder: fzflib.New().Normalize(false).Forward(true),
	}

	menu.initKeyBinder()

	return menu, nil
}

func (menu *MenuScreen) Start() *MenuScreen {

	defer func() {
		if r := recover(); r != nil {
			println(r)
			debug.PrintStack()
		}
		menu.Fini()
	}()

	screen := menu.screen
	menu.shutdownCtrl = make(chan struct{})

	for {

		if menu.isShutdown() {
			return menu
		}

		menu.refreshScreen()

		switch event := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if fn := menu.keyBinder.find(event.Key()); fn != nil {
				fn(event)
			}
		}

	}

}

func (menu *MenuScreen) Fini() {
	if !menu.finished {
		menu.screen.Fini()
		menu.finished = true
	}
}

func (menu *MenuScreen) shutdown() {
	if !menu.isShutdown() {
		close(menu.shutdownCtrl)
	}
}

func (menu *MenuScreen) isShutdown() bool {
	select {
	case <-menu.shutdownCtrl:
		return true
	default:
		return false
	}
}
