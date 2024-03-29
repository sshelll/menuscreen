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

import "github.com/gdamore/tcell/v2"

var (
	defaultContentStyle = tcell.StyleDefault.
				Background(tcell.ColorReset).
				Foreground(tcell.ColorReset)

	defaultTitleStyle = defaultContentStyle.
				Bold(true).
				Italic(true)

	defaultChosenLineStyle = tcell.StyleDefault.
				Foreground(tcell.ColorYellow).
				Background(tcell.ColorReset).
				Bold(true)

	defaultChosenLineStyleLight = tcell.StyleDefault.
					Foreground(tcell.ColorYellow).
					Background(tcell.ColorReset).
					Bold(true)

	defaultCursorColStyle = defaultContentStyle

	defaultCursorColStyleLight = defaultContentStyle

	defaultQueryStyle = defaultContentStyle.
				Italic(true)

	defaultHighlightStyle = defaultContentStyle.Bold(true).Reverse(true)
)

func SetTitleStyle(style tcell.Style) {
	defaultTitleStyle = style
}

func SetContentStyle(style tcell.Style) {
	defaultContentStyle = style
}

func SetChosenLineStyle(style tcell.Style) {
	defaultChosenLineStyle = style
}

func SetCursorColStyle(style tcell.Style) {
	defaultCursorColStyle = style
}

func SetQueryStyle(style tcell.Style) {
	defaultQueryStyle = style
}

func SetDefaultLightStyle() {
	defaultChosenLineStyle = defaultChosenLineStyleLight
	defaultCursorColStyle = defaultCursorColStyleLight
}
