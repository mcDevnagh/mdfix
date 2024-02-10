package fixers

import (
	"unicode"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Whitespace struct{}

func (f *Whitespace) Fix(parser parser.Parser, source []byte) ([]byte, error) {
	node := parser.Parse(text.NewReader(source))
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Type() == ast.TypeBlock && n.Kind() == ast.KindParagraph {
			lines := n.Lines()
			new_lines := text.NewSegments()
			for i := 0; i < lines.Len(); i++ {
				f.cut(lines.At(i), source)
			}

			n.SetLines(new_lines)
		}

		return ast.WalkContinue, nil
	})

	if source[len(source)-1] != '\n' {
		source = append(source, '\n')
	}

	return source, nil
}

func (f *Whitespace) cut(line text.Segment, source []byte) {
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
			f.cut(text.NewSegment(stop, line.Stop), source)
		}
	}
}
