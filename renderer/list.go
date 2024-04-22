// This file is part of gmnhg.

// gmnhg is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// gmnhg is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with gmnhg. If not, see <https://www.gnu.org/licenses/>.

package renderer

import (
	"bytes"
	"fmt"
	"io"
	"github.com/gomarkdown/markdown/ast"
)

var (
	itemIndent = []byte{'\t'}
	itemPrefix = []byte("* ")
)

func (r Renderer) list(w io.Writer, node *ast.List, level int) {
	// the text/gemini spec included with the current Gemini spec does
	// not specify anything about the formatting of lists of level >= 2,
	// as of now this will just render them like in Markdown
	isNumbered := (node.ListFlags & ast.ListTypeOrdered) != 0
	for number, item := range node.Children {
		item, ok := item.(*ast.ListItem)
		if !ok {
			return
		}
		isTerm := (item.ListFlags & ast.ListTypeTerm) == ast.ListTypeTerm
		if l := len(item.Children); l >= 1 {
			// add extra line break to split up definitions
			if isTerm && number > 0 {
				w.Write(lineBreak)
			}
			indent := bytes.Repeat(itemIndent, level)
			w.Write(indent)
			if isNumbered {
				prefix := fmt.Sprintf("%d. ", number+1)
				w.Write([]byte(prefix))
				replacement := []byte{}
				replacement = append(replacement, lineBreak...)
				replacement = append(replacement, indent...)
				replacement = append(replacement, bytes.Repeat(space, len(prefix))...)
				w.Write(textWithNewlineReplacement(item, replacement, true))
			} else if !isTerm {
				w.Write(itemPrefix)
				replacement := []byte{}
				replacement = append(replacement, lineBreak...)
				replacement = append(replacement, indent...)
				replacement = append(replacement, bytes.Repeat(space, len(itemPrefix))...)
				w.Write(textWithNewlineReplacement(item, replacement, true))
			}
			w.Write(lineBreak)
			if l >= 2 {
				if list, ok := item.Children[1].(*ast.List); ok {
					r.list(w, list, level+1)
				}
			}
		}
	}
}
