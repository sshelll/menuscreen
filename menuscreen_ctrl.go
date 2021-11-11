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
	"log"
	"runtime/debug"
)

func (menu *MenuScreen) Start() *MenuScreen {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("panicked: %v\n", r)
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
