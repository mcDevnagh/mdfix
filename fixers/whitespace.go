package fixers

import (
	"fmt"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Whitespace struct {
	TextWidth int
}

func (f *Whitespace) Fix(node ast.Node, source []byte, _ goldmark.Markdown) ([]byte, error) {
	if f.TextWidth <= 0 {
		return nil, fmt.Errorf("invalid TextWidth (%d). Value must be greater than 0", f.TextWidth)
	}

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
	target_len := line.Start + f.TextWidth
	if line.Stop > target_len {
		var stop int
		bracket_stack := 0
		parenthesis_stack := 0
		for stop = line.Stop - 1; stop > target_len; stop-- {
			stack(&bracket_stack, &parenthesis_stack, source[stop])
		}

		for stop = target_len; stop > line.Start; stop-- {
			if stack(&bracket_stack, &parenthesis_stack, source[stop]) {
				continue
			} else if unicode.IsSpace(rune(source[stop])) && bracket_stack == 0 && parenthesis_stack == 0 {
				source[stop] = '\n'
				break
			}
		}

		if stop != line.Start {
			f.cut(text.NewSegment(stop, line.Stop), source)
		}
	}
}

func stack(bracket_stack, parenthesis_stack *int, char byte) bool {
	if char == ']' {
		*bracket_stack++
		return true
	} else if char == '[' {
		*bracket_stack--
		if *bracket_stack < 0 {
			*bracket_stack = 0
		}
		return true
	} else if char == ')' {
		*parenthesis_stack++
		return true
	} else if char == '(' {
		*parenthesis_stack--
		if *parenthesis_stack < 0 {
			*parenthesis_stack = 0
		}
		return true
	} else {
		return false
	}
}
