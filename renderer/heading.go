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
	"io"
	"github.com/gomarkdown/markdown/ast"
)

var majorUnderline = []byte{'='}
var minorUnderline = []byte{'-'}

func (r Renderer) heading(w io.Writer, node *ast.Heading, entering bool) {
	if !entering {
		text := ast.GetFirstChild(node).AsLeaf()
		w.Write(text.Literal)
		w.Write(lineBreak)
		var s []byte;
		if node.Level == 1 {
			s = majorUnderline
		} else {
			s = minorUnderline
		}
		for i := 0; i < len(text.Literal); i++ {
			w.Write(s)
		}
		w.Write(lineBreak)
	}
}
