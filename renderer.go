// This is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as
// published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

// This is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty
// of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

// You should have received a copy of the GNU General Public License along with this. If not, see
// <https://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"fmt"
	"github.com/gomarkdown/markdown/ast"
	"io"
	"strings"
	"unicode/utf8"
)

var (
	newline        = "\n"
	space          = " "
	rule           = "-----------------------------------------------------------------"
	majorUnderline = "="
	minorUnderline = "-"
)

// Wrapper is a wordwrapper based on https://godoc.org/github.com/karrick/golinewrap
type Wrapper struct {
	line         *bytes.Buffer
	max          int // max number of runes to fill for each line
	remaining    int // remaining runes in the line buffer
	prefixLength int
	prefix       string
}

// Renderer implements markdown.Renderer based on https://github.com/tdemin/gmnhg
type Renderer struct {
	buf *Wrapper
}

// NewRenderer returns a new Renderer.
func NewRenderer() Renderer {
	wrapper := Wrapper{
		line:      bytes.NewBuffer(make([]byte, 0, 73)),
		max:       72,
		remaining: 72,
		prefix:    "",
	}
	return Renderer{buf: &wrapper}
}

// setPrefix changes the prefix.
func (w *Wrapper) setPrefix(prefix string) {
	w.prefix = prefix
	w.prefixLength = utf8.RuneCountInString(prefix)
}

// write writes a string to the buffer, without adding a space or a newline.
func (buf *Wrapper) write(w io.Writer, s string) {
	required := utf8.RuneCountInString(s) + 1
	if buf.remaining < required {
		buf.newline(w)
	}
	buf.line.WriteString(s)
	buf.remaining -= (required - 1)
}

// writePrefix writes the prefix, at the beginning of a line, if any.
func (buf *Wrapper) writePrefix() {
	if buf.max == buf.remaining && buf.prefix != "" {
		buf.line.WriteString(buf.prefix)
		buf.remaining -= buf.prefixLength
	}
}

// writeWords appends the words in the text to the current paragraph. The prefix must be set beforehand, if at all. The
// text ends with a space. Calling newline will remove that space again. If you don't want this space, use write.
// Calling newline also flushes the buffer.
func (buf *Wrapper) writeWords(w io.Writer, text string) {
	for _, word := range strings.Fields(text) {
		buf.writePrefix()
		buf.write(w, word+space)
	}
}

// newline appends newline to line buffer then flushes to underlying writer because this library is line based.
func (buf *Wrapper) newline(w io.Writer) {
	b := buf.line.Bytes()
	if l := len(b); l > 0 && b[l-1] == ' ' {
		// remove final space character from line buffer
		buf.line.Truncate(l - 1)
	}
	buf.line.WriteString(newline)
	buf.remaining = buf.max
	buf.flush(w)
}

// flush writes any remaining runes in the line out to the writer, without adding a newline.
func (buf *Wrapper) flush(w io.Writer) {
	buf.line.WriteTo(w)
}

// RenderHeader implements Renderer.RenderHeader(). As there is no header, there is nothing to do, here.
func (r Renderer) RenderHeader(w io.Writer, node ast.Node) {}

// RenderFooter implements Renderer.RenderFooter(). As there is no footer, there is nothing to do, here.
func (r Renderer) RenderFooter(w io.Writer, node ast.Node) {
	r.buf.newline(w)
}

// RenderNode implements Renderer.RenderNode(). This goes through every node and fills a buffer with words. As soon as
// there are enough words for a line, it is written to the Writer.
func (r Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.BlockQuote:
		r.buf.setPrefix("> ")
	case *ast.HorizontalRule:
		r.buf.newline(w)
		r.buf.write(w, rule)
		r.buf.newline(w)
		r.buf.newline(w)
	case *ast.Heading:
		if !entering {
			r.buf.newline(w)
			text := ast.GetFirstChild(node).AsLeaf()
			var s string
			if node.Level == 1 {
				s = majorUnderline
			} else {
				s = minorUnderline
			}
			r.buf.write(w, strings.Repeat(s, len(text.Literal)))
			r.buf.newline(w)
			r.buf.newline(w)
		}
	case *ast.Paragraph:
		r.buf.setPrefix("")
	case *ast.Text:
		r.buf.writeWords(w, string(node.Literal))
	case *ast.Document:
	default:
		text := node.AsLeaf()
		if text != nil {
			r.buf.writeWords(w, string(text.Literal))
		} else {
			fmt.Printf("Ignoring %T\n", node)
		}
	}
	return ast.GoToNext
}
