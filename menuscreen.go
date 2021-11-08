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
)

// MenuScreen is a visible selector for input.
//
// Terminal will show the content below in the normal mode:
//  -----------------
// | ${Title}:       |
// | ▸ ${1st line}   |
// |   ${2nd line}   |
// |   ${3rd line}   |
//  -----------------
//
// Once slash was pressed and MenuScreen entered query mode,
// the screen will be like:
//  -----------------------
// | ${Title}:             |
// | /query                |
// | ▸ ${1st matched line} |
// |   ${2nd matched line} |
// |   ${3rd matched line} |
//  -----------------------
type MenuScreen struct {
	screen       tcell.Screen
	keyBinder    *keyBinder
	shutdownCtrl chan struct{}
	mode         screenMode
	cursorY      int
	input        string
	title        string
	lines        []string
	matchedLns   matchedLines
	confirmed    bool
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
		screen:     screen,
		lines:      make([]string, 0, 16),
		matchedLns: make([]*matchedLine, 0, 16),
		mode:       modeN,
		cursorY:    0,
		input:      "",
		title:      "Menu",
	}

	menu.initKeyBinder()

	return menu, nil
}
