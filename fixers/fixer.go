package fixers

import (
	"github.com/yuin/goldmark/parser"
)

type Fixer interface {
	Fix(parser parser.Parser, source []byte) (fixed []byte, err error)
}
