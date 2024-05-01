package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestText(t *testing.T) {
	body := []byte("This is text.\n")
	expected := "This is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestEmphasis(t *testing.T) {
	body := []byte("This is *italic text*.\n")
	expected := "This is italic text.\n"
	assert.Equal(t, expected, render(body))
}

func TestStrong(t *testing.T) {
	body := []byte("This is **bold text**.\n")
	expected := "This is bold text.\n"
	assert.Equal(t, expected, render(body))
}

func TestInlineCode(t *testing.T) {
	body := []byte("This is `some code`.\n")
	expected := "This is some code.\n"
	assert.Equal(t, expected, render(body))
}

func TestLink(t *testing.T) {
	body := []byte("This is [a link](https://example.com).\n")
	expected := "This is a link.\n"
	assert.Equal(t, expected, render(body))
}

func TestTwoLines(t *testing.T) {
	body := []byte(`This is text.
And this is text.`)
	expected := "This is text. And this is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestLongLine(t *testing.T) {
	body := append([]byte(strings.Repeat("This is text. ", 10)), []byte("The End.\n")...)
	expected := `This is text. This is text. This is text. This is text. This is text.
This is text. This is text. This is text. This is text. This is text.
The End.
`
	assert.Equal(t, expected, render(body))
}

func TestHeading1(t *testing.T) {
	body := []byte(`# Heading

This is text.
`)
	expected := `Heading
=======

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestHeading2(t *testing.T) {
	body := []byte(`## Heading

This is text.
`)
	expected := `Heading
-------

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestRule(t *testing.T) {
	body := []byte(`This is text.

----

This is text.`)
	expected := `This is text.

----------------------------------------------------------------------

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestQuote(t *testing.T) {
	body := []byte(`This is text.

> Hier kommt ein Zitat. Es ist ein langes Zitat. Mit mehreren Sätzen. Und mehreren Zeilen.

This is text.`)
	expected := `This is text.

> Hier kommt ein Zitat. Es ist ein langes Zitat. Mit mehreren Sätzen.
> Und mehreren Zeilen.

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestQuote2(t *testing.T) {
	// strange but true: the empty line between two quoted lines does not result in two quoted blocks!
	body := []byte(`This is text.

> Hier kommt ein Zitat.

> Und das ist ein anderes Zitat.

This is text.`)
	expected := `This is text.

> Hier kommt ein Zitat.
>
> Und das ist ein anderes Zitat.

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestQuote3(t *testing.T) {
	body := []byte(`This is text.

> Hier kommt ein Zitat.
>
> Mit zwei Paragraphen.

This is text.`)
	expected := `This is text.

> Hier kommt ein Zitat.
>
> Mit zwei Paragraphen.

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestCode(t *testing.T) {
	body := []byte("This is text.\n\n```\nsome code\nsome more code\n```\n\nThis is text.")
	expected := `This is text.

    some code
    some more code

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestUnorderedList(t *testing.T) {
	body := []byte(`This is text.

* This is an item.

This is text.`)
	expected := `This is text.

* This is an item.

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestUnorderedList2(t *testing.T) {
	body := []byte(`This is text.

* This is a very long item. This is a very long item. This is a very long item. This is a very long item.

This is text.`)
	expected := `This is text.

* This is a very long item. This is a very long item. This is a very
  long item. This is a very long item.

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestUnorderedList3(t *testing.T) {
	body := []byte(`This is text.

* This is an item.

* This is a very long item. This is a very long item. This is a very long item. This is a very long item.

This is text.`)
	expected := `This is text.

* This is an item.

* This is a very long item. This is a very long item. This is a very
  long item. This is a very long item.

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestUnorderedList4(t *testing.T) {
	body := []byte(`This is text.

* This is an item.

    * This is a very long item. This is a very long item. This is a very long item. This is a very long item.

This is text.`)
	expected := `This is text.

* This is an item.

  * This is a very long item. This is a very long item. This is a very
    long item. This is a very long item.

This is text.
`
	assert.Equal(t, expected, render(body))
}


func TestOrderedList(t *testing.T) {
	body := []byte(`This is text.

1. This is an item.

This is text.`)
	expected := `This is text.

1. This is an item.

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestOrderedList2(t *testing.T) {
	body := []byte(`This is text.

1. This is a very long item. This is a very long item. This is a very long item. This is a very long item.

This is text.`)
	expected := `This is text.

1. This is a very long item. This is a very long item. This is a very
   long item. This is a very long item.

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestOrderedList3(t *testing.T) {
	body := []byte(`This is text.

1. This is an item.

1. This is a very long item. This is a very long item. This is a very long item. This is a very long item.

This is text.`)
	expected := `This is text.

1. This is an item.

2. This is a very long item. This is a very long item. This is a very
   long item. This is a very long item.

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestOrderedList4(t *testing.T) {
	body := []byte(`This is text.

1. This is an item.

    1. This is a very long item. This is a very long item. This is a very long item. This is a very long item.

This is text.`)
	expected := `This is text.

1. This is an item.

  1. This is a very long item. This is a very long item. This is a
     very long item. This is a very long item.

This is text.
`
	assert.Equal(t, expected, render(body))
}

func TestTable(t *testing.T) {
	body := []byte(`This is text.

Name    | Age
--------|------
Bob     | 27
Alice   | 23

This is text.`)
	expected := `This is text.

+-------+-----+
| NAME  | AGE |
+-------+-----+
| Bob   |  27 |
| Alice |  23 |
+-------+-----+

This is text.
`
	assert.Equal(t, expected, render(body))
}

func render(input []byte) string {
	p := parser.New()
	ast := p.Parse(input)
	content := markdown.Render(ast, NewRenderer())
	return string(content)
}
