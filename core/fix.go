package core

import (
	"io"
	"unicode"

	"github.com/Kunde21/markdownfmt/v3/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func Fix(source []byte, dest io.Writer) error {
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

	node = parser.Parse(text.NewReader(source))
	return md.Renderer().Render(dest, source, node)
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
