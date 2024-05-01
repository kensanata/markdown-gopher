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
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/olekukonko/tablewriter"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	newline        = "\n"
	space          = " "
	rule           = "----------------------------------------------------------------------"
	bullet         = "*"
	majorUnderline = "="
	minorUnderline = "-"
)

// Counter is a counter for list items.
type Counter struct {
	counter []int
}

// Wrapper is a wordwrapper based on the ideas in https://godoc.org/github.com/karrick/golinewrap.
type Wrapper struct {
	Counter
	first        bool // is this the first block of the document
	line         *bytes.Buffer
	max          int // max number of runes to fill for each line
	remaining    int // remaining runes in the line buffer
	prefixLength int // how long the prefix is
	prefix       string // the prefix for lines
	prefixNext   string // the next prefix (for list items)
	prefixSkip   bool // skip prefix when prefix changes
	tab          *tablewriter.Table
	header       bool // is this a header row for the table
	footer       bool // is this a footer row for the table
	row          []string
}

// Renderer implements markdown. The initial idea of how it was going to work are on https://github.com/tdemin/gmnhg.
type Renderer struct {
	buf         *Wrapper
}

// push adds a new counter starting with 0
func (c *Counter) push() {
	c.counter = append(c.counter, 0)
}

// pop removes the last counter
func (c *Counter) pop() {
	c.counter = c.counter[:len(c.counter)-1]
}

// inc increases the current counter by one
func (c *Counter) inc() {
	c.counter[len(c.counter)-1]++
}

// val returns the current counter value
func (c *Counter) value() int {
	return c.counter[len(c.counter)-1]
}

// trim removes the trailing space of the buffer, if any.
func (buf *Wrapper) trim() {
	b := buf.line.Bytes()
	if l := len(b); l > 0 && b[l-1] == ' ' {
		buf.line.Truncate(l - 1)
	}
}

// setPrefix changes the prefix.
func (buf *Wrapper) setPrefix(prefix string) {
	buf.prefix = prefix
	buf.prefixLength = utf8.RuneCountInString(prefix)
	if (buf.prefix == "") {
		buf.prefixNext = ""
	}
}

// writePrefix writes the prefix, but only at the beginning of the line (and if a prefix is set). Call it anytime
// something is written to the line.
func (buf *Wrapper) writePrefix() {
	if buf.max == buf.remaining && buf.prefix != "" {
		buf.line.WriteString(buf.prefix)
		buf.remaining -= buf.prefixLength
		if buf.prefixNext != "" {
			buf.setPrefix(buf.prefixNext)
			buf.prefixNext = ""
		}
	}
}

// write writes a string to the buffer. If necessary, a newline is prepended and the previous line is flushed to the
// writer. No space or a newline is added at the end. Use this for a horizontal rule or to underline headings.
func (buf *Wrapper) write(w io.Writer, s string) {
	required := utf8.RuneCountInString(s) + 1
	if buf.remaining < required {
		buf.newline(w)
	}
	buf.writePrefix()
	buf.line.WriteString(s)
	buf.remaining -= (required - 1)
}

// writeWords appends the words in the text to the line. The prefix must already be set, if at all. If the text starts
// or ends with whitespace, a space is prefixed or suffixed, respectively. The words are written using the write
// function, which might call newline, which flushes the line. At that point, a trailing space is going to be trimmed
// from the line.
func (buf *Wrapper) writeWords(w io.Writer, text string) {
	// if the text starts with whitespace, prepend a single space
	rune, size := utf8.DecodeRuneInString(text)
	if size > 0 && unicode.IsSpace(rune) {
		buf.write(w, space)
	}
	// always add a single space after every word
	for _, word := range strings.Fields(text) {
		buf.write(w, word+space)
	}
	// if the text doesn't end with whitespace, trim that last space again
	rune, size = utf8.DecodeLastRuneInString(text)
	if size == 0 || !unicode.IsSpace(rune) {
		buf.trim()
	}
}

// newline trims a trailing space from the line, appends a newline and flushes the line to the underlying writer. Every
// block element must end with a call to newline or the last line of the document will not be flushed.
func (buf *Wrapper) newline(w io.Writer) {
	buf.trim()
	buf.line.WriteString(newline)
	buf.remaining = buf.max
	buf.line.WriteTo(w)
}

// NewRenderer returns a new Renderer.
func NewRenderer() Renderer {
	wrapper := Wrapper{
		first:      true,
		line:       bytes.NewBuffer(make([]byte, 0, 73)),
		max:        72,
		remaining:  72,
		prefix:     "",
		prefixNext: "",
		prefixSkip: false,
	}
	return Renderer{buf: &wrapper}
}

// RenderHeader implements Renderer.RenderHeader(). As there is no header, there is nothing to do, here.
func (r Renderer) RenderHeader(w io.Writer, node ast.Node) {}

// RenderFooter implements Renderer.RenderFooter(). As there is no footer, there is nothing to do, here.
func (r Renderer) RenderFooter(w io.Writer, node ast.Node) {}

// RenderNode implements Renderer.RenderNode(). This goes through every node and fills a buffer with words. As soon as
// there are enough words for a line, it is written to the Writer.
func (r Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	fmt.Printf("%T %v\n", node, entering)
	switch node := node.(type) {
	case *ast.BlockQuote:
		if entering {
			r.buf.setPrefix("> ")
			r.buf.prefixSkip = true
		} else {
			r.buf.setPrefix("")
		}
	case *ast.List:
		if entering {
			r.buf.push()
		} else {
			r.buf.pop()
		}
	case *ast.ListItem:
		if entering {
			r.buf.inc()
			isOrdered := (node.ListFlags & ast.ListTypeOrdered) == ast.ListTypeOrdered
			isDefinition := (node.ListFlags & ast.ListTypeTerm) == ast.ListTypeTerm
			isUnordered := !isOrdered && !isDefinition
			indentation := strings.Repeat(space, 2 * (len(r.buf.counter)-1))
			if isUnordered {
				r.buf.setPrefix(indentation + bullet + space)
				r.buf.prefixNext = strings.Repeat(space, len(r.buf.prefix))
			} else if isOrdered {
				r.buf.setPrefix(fmt.Sprintf("%s%d. ", indentation, r.buf.value()))
				r.buf.prefixNext = strings.Repeat(space, len(r.buf.prefix))
			} else if r.buf.value() > 1 {
				// definition lists get extra line breaks
				r.buf.newline(w)
				r.buf.setPrefix(indentation)
			}
			r.buf.prefixSkip = true
		} else {
			r.buf.setPrefix("")
		}
	case *ast.HorizontalRule:
		r.paragraphSeparator(w)
		r.buf.write(w, rule)
		r.buf.newline(w)
	case *ast.Heading:
		if entering {
			r.paragraphSeparator(w)
		} else {
			// After the text of the heading, add underlining
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
		}
	case *ast.Paragraph:
		if entering {
			r.paragraphSeparator(w)
		} else {
			r.buf.newline(w)
		}
	case *ast.CodeBlock:
		newline := []byte("\n")
		separator := []byte("\n    ")
		w.Write(separator)
		w.Write(bytes.Replace(node.Literal, newline,  separator, bytes.Count(node.Literal, newline)-1))
	case *ast.Table:
		// The tablewriter needs to be fed cells in rows.
		if entering {
			r.buf.tab = tablewriter.NewWriter(w)
			r.buf.newline(w)
		} else {
			r.buf.tab.Render()
		}
	case *ast.TableRow:
		if entering {
			r.buf.row = make([]string, 0)
		} else if r.buf.header {
			r.buf.tab.SetHeader(r.buf.row)
		} else if r.buf.footer {
			r.buf.tab.SetFooter(r.buf.row)
		} else {
			r.buf.tab.Append(r.buf.row)
		}
	case *ast.TableHeader:
		r.buf.header = entering
	case *ast.TableFooter:
		r.buf.footer = entering
	case *ast.TableBody:
	case *ast.TableCell:
		if entering {
			// render the children of the table cell (without the table cell itself)
			doc := &ast.Document{}
			doc.SetChildren(node.GetChildren())
			s := string(markdown.Render(doc, NewRenderer()))
			r.buf.row = append(r.buf.row, s)
			return ast.SkipChildren
		}
	case *ast.Text:
		r.buf.writeWords(w, string(node.Literal))
	case *ast.Code:
		r.buf.writeWords(w, string(node.Literal))
	case *ast.Emph:
	case *ast.Strong:
	case *ast.Document:
		if !entering {
			r.buf.line.WriteTo(w) // flush for TableCell
		}
	default:
		text := node.AsLeaf()
		if text != nil {
			r.buf.writeWords(w, string(text.Literal))
		}
	}
	return ast.GoToNext
}

// paragraphSeparator writes a paragraph separator unless this is the first paragraph. Furthermore, when a new type
// starts (such as a quote), the prefix is skipped once before the newline. The result is that when a quoted paragraph
// follows a regular paragraph, the line is empty but if a quoted paragraph follows another quoted paragraph, the line
// is quoted.
func (r Renderer) paragraphSeparator(w io.Writer) {
	if r.buf.first {
		r.buf.first = false
	} else {
		if r.buf.prefixSkip {
			r.buf.prefixSkip = false
		} else {
			r.buf.writePrefix()
		}
		r.buf.newline(w)
	}
}
