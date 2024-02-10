package core

import (
	"io"

	"github.com/Kunde21/markdownfmt/v3/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"gitlab.com/mcdonagh/mdfix/fixers"
)

func Fix(source []byte, dest io.Writer) error {
	return fix(source, dest, []fixers.Fixer{
		&fixers.Whitespace{},
	})
}

func fix(source []byte, dest io.Writer, fixers []fixers.Fixer) error {
	mr := markdown.NewRenderer()
	mr.AddMarkdownOptions(markdown.WithSoftWraps())
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAttribute()),
		goldmark.WithRenderer(mr),
	)

	parser := md.Parser()
	for _, fixer := range fixers {
		var err error
		source, err = fixer.Fix(parser, source)
		if err != nil {
			return err
		}
	}

	node := parser.Parse(text.NewReader(source))
	return md.Renderer().Render(dest, source, node)
}
