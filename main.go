package main

import (
	"fmt"
	"git.mills.io/prologic/go-gopher"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"io"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"src.alexschroeder.ch/markdown-gopher/renderer"
	"strings"
)

func main() {
	gopher.HandleFunc("/", serve)
	port := 7000
	addr := fmt.Sprintf("localhost:%d", port)
	fmt.Printf("Listening on %s\n", addr)
	log.Fatal(gopher.ListenAndServe(addr, nil))
}

func serve(w gopher.ResponseWriter, r *gopher.Request) {
	const (
		unknown = iota
		file
		page
		dir
	)
	t := unknown
	path := r.Selector
	fp := filepath.Join(".", filepath.FromSlash(path))
	fmt.Println("Path: " + fp)
	fi, err := os.Stat(fp + ".md")
	if err == nil {
		if fi.IsDir() {
			t = dir // directory ending in ".md"
		} else {
			t = page
		}
	} else {
		fi, err = os.Stat(fp)
		if err == nil {
			if fi.IsDir() {
				t = dir
			} else {
				t = file
			}
		}
	}
	// if nothing was found, abort
	if t == unknown {
		fmt.Fprint(w, "no info available\r\n")
		return
	}
	// directories are redirected to the index page
	if t == dir {
		menu(w, r, fp)
		return
	}
	// if the file exists, serve it
	if t == file {
		file, err := os.Open(fp)
		if err != nil {
			fmt.Fprint(w, "unable to open file\r\n")
			log.Println(err)
			return
		}
		// copy file
		_, err = io.Copy(w, file)
		if err != nil {
			fmt.Fprint(w, "unable to copy file\r\n")
			log.Println(err)
			return
		}
		return
	}
	md, err := load(fp)
	if err != nil {
		fmt.Fprint(w, "unable to load file\r\n")
		log.Println(err)
		return
	}
	w.Write(md)
}

func menu(w gopher.ResponseWriter, r *gopher.Request, path string) {
	fi, err := os.ReadFile(filepath.Join(path, "index.md"))
	if err != nil {
		filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if p == path {
				return nil;
			} else if info.IsDir() {
				return filepath.SkipDir
			} else if strings.HasSuffix(path, ".md") {
				name := path[:len(path)-3]
				fmt.Fprintf(w, "0%s\t%s\t%s\t%d\r\n", name, name, r.LocalHost, r.LocalPort)
			}
			return nil
		})
	} else {
		re := regexp.MustCompile(`(?m)^\* \[(.*?)\]\((.*?)\)`)
		for _, m := range re.FindAllSubmatch(fi, -1) {
			fmt.Fprintf(w, "0%s\t%s\t%s\t%d\r\n", string(m[1]), filepath.Join(path, string(m[2])), r.LocalHost, r.LocalPort)
		}
	}
}

func load(path string) ([]byte, error) {
	md, err := os.ReadFile(path + ".md")
	if err != nil {
		return nil, err
	}
	ast := markdown.Parse(md, wikiParser())
	content := markdown.Render(ast, renderer.NewRenderer())
	return content, nil
}

// wikiParser returns a parser with the Oddmu specific changes.
// Specifically: [[wiki links]], #hash_tags, @webfinger@accounts.
// It also uses the CommonExtensions without MathJax ($).
func wikiParser() *parser.Parser {
	extensions := parser.CommonExtensions & ^parser.MathJax
	parser := parser.NewWithExtensions(extensions)
	prev := parser.RegisterInline('[', nil)
	parser.RegisterInline('[', wikiLink(prev))
	parser.RegisterInline('#', hashtag)
	return parser
}

// wikiLink returns an inline parser function. This indirection is
// required because we want to call the previous definition in case
// this is not a wikiLink.
func wikiLink(fn func(p *parser.Parser, data []byte, offset int) (int, ast.Node)) func(p *parser.Parser, data []byte, offset int) (int, ast.Node) {
	return func(p *parser.Parser, original []byte, offset int) (int, ast.Node) {
		data := original[offset:]
		n := len(data)
		// minimum: [[X]]
		if n < 5 || data[1] != '[' {
			return fn(p, original, offset)
		}
		i := 2
		for i+1 < n && data[i] != ']' && data[i+1] != ']' {
			i++
		}
		text := data[2 : i+1]
		link := &ast.Link{
			Destination: []byte(url.PathEscape(string(text))),
		}
		ast.AppendChild(link, &ast.Text{Leaf: ast.Leaf{Literal: text}})
		return i + 3, link
	}
}

// hashtag turns hashtags into plain text by prefixing them with a zero-width space.
func hashtag(p *parser.Parser, data []byte, offset int) (int, ast.Node) {
	data = data[offset:]
	i := 0
	n := len(data)
	for i < n && !parser.IsSpace(data[i]) {
		i++
	}
	if i == 0 {
		return 0, nil
	}
	text := append([]byte("\u200b"), data...)
	return i, &ast.Text{Leaf: ast.Leaf{Literal: text}}
}
