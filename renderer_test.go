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

func TestTwoLines(t *testing.T) {
	body := []byte("This is text.\nAnd this is text.")
	expected := "This is text. And this is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestLongLine(t *testing.T) {
	body := append([]byte(strings.Repeat("This is text. ", 10)), []byte("The End.\n")...)
	expected := "This is text. This is text. This is text. This is text. This is text.\n" +
		"This is text. This is text. This is text. This is text. This is text.\n" +
		"The End.\n"
	assert.Equal(t, expected, render(body))
}

func TestHeading1(t *testing.T) {
	body := []byte("# Heading\n\nThis is text.\n")
	expected := "Heading\n=======\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestHeading2(t *testing.T) {
	body := []byte("## Heading\n\nThis is text.\n")
	expected := "Heading\n-------\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestRule(t *testing.T) {
	body := []byte("This is text.\n\n----\n\nThis is text.")
	expected := "This is text.\n\n----------------------------------------------------------------------\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestQuote(t *testing.T) {
	body := []byte("This is text.\n\n> Hier kommt ein Zitat. Es ist ein langes Zitat.\n" +
		"Mit mehreren Sätzen. Und mehreren Zeilen.\n\nThis is text.")
	expected := "This is text.\n\n> Hier kommt ein Zitat. Es ist ein langes Zitat. Mit mehreren Sätzen.\n> Und mehreren Zeilen.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestQuote2(t *testing.T) {
	// strange but true: the empty line between two quoted lines does not result in two quoted blocks!
	body := []byte("This is text.\n\n> Hier kommt ein Zitat.\n\n> Und das ist ein anderes Zitat.\n\nThis is text.")
	expected := "This is text.\n\n> Hier kommt ein Zitat.\n>\n> Und das ist ein anderes Zitat.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestQuote3(t *testing.T) {
	body := []byte("This is text.\n\n> Hier kommt ein Zitat.\n>\n> Mit zwei Paragraphen.\n\nThis is text.")
	expected := "This is text.\n\n> Hier kommt ein Zitat.\n>\n> Mit zwei Paragraphen.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestUnorderedList(t *testing.T) {
	body := []byte("This is text.\n\n* This is an item.\n\nThis is text.")
	expected := "This is text.\n\n* This is an item.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestUnorderedList2(t *testing.T) {
	body := []byte("This is text.\n\n* This is a very long item. This is a very long item. This is a very long item. This is a very long item.\n\nThis is text.")
	expected := "This is text.\n\n* This is a very long item. This is a very long item. This is a very\n  long item. This is a very long item.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestUnorderedList3(t *testing.T) {
	body := []byte("This is text.\n\n* This is an item.\n\n* This is a very long item. This is a very long item. This is a very long item. This is a very long item.\n\nThis is text.")
	expected := "This is text.\n\n* This is an item.\n\n* This is a very long item. This is a very long item. This is a very\n  long item. This is a very long item.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestUnorderedList4(t *testing.T) {
	body := []byte("This is text.\n\n* This is an item.\n\n    * This is a very long item. This is a very long item. This is a very long item. This is a very long item.\n\nThis is text.")
	expected := "This is text.\n\n* This is an item.\n\n  * This is a very long item. This is a very long item. This is a very\n    long item. This is a very long item.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}


func TestOrderedList(t *testing.T) {
	body := []byte("This is text.\n\n1. This is an item.\n\nThis is text.")
	expected := "This is text.\n\n1. This is an item.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestOrderedList2(t *testing.T) {
	body := []byte("This is text.\n\n1. This is a very long item. This is a very long item. This is a very long item. This is a very long item.\n\nThis is text.")
	expected := "This is text.\n\n1. This is a very long item. This is a very long item. This is a very\n   long item. This is a very long item.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestOrderedList3(t *testing.T) {
	body := []byte("This is text.\n\n1. This is an item.\n\n1. This is a very long item. This is a very long item. This is a very long item. This is a very long item.\n\nThis is text.")
	expected := "This is text.\n\n1. This is an item.\n\n2. This is a very long item. This is a very long item. This is a very\n   long item. This is a very long item.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func TestOrderedList4(t *testing.T) {
	body := []byte("This is text.\n\n1. This is an item.\n\n    1. This is a very long item. This is a very long item. This is a very long item. This is a very long item.\n\nThis is text.")
	expected := "This is text.\n\n1. This is an item.\n\n  1. This is a very long item. This is a very long item. This is a\n     very long item. This is a very long item.\n\nThis is text.\n"
	assert.Equal(t, expected, render(body))
}

func render(input []byte) string {
	p := parser.New()
	ast := p.Parse(input)
	content := markdown.Render(ast, NewRenderer())
	return string(content)
}
