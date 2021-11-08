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
)

type keyBinder struct {
	mapping map[tcell.Key]func(*tcell.EventKey)
}

func (kb *keyBinder) bind(fn func(*tcell.EventKey), keys ...tcell.Key) *keyBinder {

	if kb.mapping == nil {
		kb.mapping = make(map[tcell.Key]func(*tcell.EventKey))
	}

	if len(keys) == 0 {
		return kb
	}

	for _, k := range keys {
		if _, ok := kb.mapping[k]; ok {
			// do not try to recover this
			panic(fmt.Sprintf("key [%v] bind twice", k))
		}
		kb.mapping[k] = fn
	}

	return kb

}

func (kb *keyBinder) find(key tcell.Key) func(*tcell.EventKey) {
	return kb.mapping[key]
}
