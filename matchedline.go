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

type matchedLines []*matchedLine

func strSliceToMatchedLines(strs []string) matchedLines {
	lns := make(matchedLines, 0, len(strs))
	for i, s := range strs {
		lns = append(lns, &matchedLine{
			idx:     i,
			content: s,
		})
	}
	return lns
}

type matchedLine struct {
	idx     int
	content string
	pos     []int
}

func (lns matchedLines) Content() []string {
	res := make([]string, 0, len(lns))
	for _, l := range lns {
		res = append(res, l.content)
	}
	return res
}
