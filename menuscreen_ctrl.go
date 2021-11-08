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

func (menu *MenuScreen) Start() {

	screen := menu.screen

	defer func() {
		if r := recover(); r != nil {
			log.Printf("panicked: %v\n", r)
			debug.PrintStack()
		}
		screen.Fini()
	}()

	menu.shutdownCtrl = make(chan struct{})

	for {

		if menu.IsShutdown() {
			return
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

func (menu *MenuScreen) Shutdown() {
	close(menu.shutdownCtrl)
}

func (menu *MenuScreen) IsShutdown() bool {
	select {
	case <-menu.shutdownCtrl:
		return true
	default:
		return false
	}
}