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

func render(input []byte) string {
	p := parser.New()
	ast := p.Parse(input)
	content := markdown.Render(ast, NewRenderer())
	return string(content)
}
