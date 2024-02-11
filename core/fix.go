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

func Fix(source []byte, dest io.Writer, options fixers.Options) error {
	_fixers := make([]fixers.Fixer, 0, 2)
	if options.TextWidth > 0 {
		_fixers = append(_fixers, &fixers.Whitespace{
			TextWidth: 80,
		})
	}

	if options.WorkDir != "" {
		_fixers = append(_fixers, &fixers.Links{
			TopDir:  options.TopDir,
			WorkDir: options.WorkDir,
		})
	}

	return fix(source, dest, _fixers)
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
		source, err = fixer.Fix(parser.Parse(text.NewReader(source)), source, md)
		if err != nil {
			return err
		}
	}

	node := parser.Parse(text.NewReader(source))
	return md.Renderer().Render(dest, source, node)
}
