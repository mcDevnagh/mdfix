package fixers

import "github.com/yuin/goldmark/ast"

type Fixer interface {
	Fix(node ast.Node, source []byte) (fixed []byte, err error)
}
