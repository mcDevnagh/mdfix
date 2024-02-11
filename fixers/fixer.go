package fixers

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
)

type Fixer interface {
	Fix(node ast.Node, source []byte, md goldmark.Markdown) (fixed []byte, err error)
}
