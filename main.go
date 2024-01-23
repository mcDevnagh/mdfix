package main

import (
	"flag"
	"io"
	"os"
	"unicode"

	"github.com/Kunde21/markdownfmt/v3/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func main() {
	in_place := flag.Bool("i", false, "inplace")
	flag.Parse()
	path := flag.Arg(0)
	var file *os.File
	var err error
	if path == "" {
		file = os.Stdin
	} else {
		file, err = os.Open(path)
		if err != nil {
			panic(err)
		}
	}

	source, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	mr := markdown.NewRenderer()
	mr.AddMarkdownOptions(markdown.WithSoftWraps())
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAttribute()),
		goldmark.WithRenderer(mr),
	)

	parser := md.Parser()
	node := parser.Parse(text.NewReader(source))
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Type() == ast.TypeBlock && n.Kind() == ast.KindParagraph {
			lines := n.Lines()
			new_lines := text.NewSegments()
			for i := 0; i < lines.Len(); i++ {
				cut(lines.At(i), source)
			}

			n.SetLines(new_lines)
		}

		return ast.WalkContinue, nil
	})

	if source[len(source)-1] != '\n' {
		source = append(source, '\n')
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	if *in_place {
		file, err = os.Create(path)
		if err != nil {
			panic(err)
		}
	} else {
		file = os.Stdout
	}

	node = parser.Parse(text.NewReader(source))
	err = md.Renderer().Render(file, source, node)
	if err != nil {
		panic(err)
	}
}

func cut(line text.Segment, source []byte) {
	target_len := line.Start + 80
	if line.Stop > target_len {
		var stop int
		for stop = target_len; stop > line.Start; stop-- {
			if unicode.IsSpace(rune(source[stop])) {
				source[stop] = '\n'
				break
			}
		}

		if stop != line.Start {
			cut(text.NewSegment(stop, line.Stop), source)
		}
	}
}
