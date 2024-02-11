package fixers

import (
	"bytes"
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
)

type Links struct {
	TopDir  string
	WorkDir string
}

func (f *Links) Fix(node ast.Node, source []byte, md goldmark.Markdown) ([]byte, error) {
	topDir := f.TopDir
	if topDir == "" {
		if runtime.GOOS == "windows" {
			topDir = os.Getenv("SystemDrive")
		} else {
			topDir = "/"
		}
	}

	var fixed bytes.Buffer
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if n.Kind() == ast.KindLink {
				link, ok := n.(*ast.Link)
				if ok {
					dest := link.Destination
					idx := bytes.IndexRune(link.Destination, '#')
					var anchor []byte
					if idx != 0 {
						if idx > 0 {
							anchor = dest[idx:]
							dest = dest[:idx]
						}

						file := string(dest)
						parent, err := getParent(file, topDir, f.WorkDir)
						if err == nil {
							file, err = filepath.Rel(f.WorkDir, path.Join(parent, file))
							if err == nil {
								dest = append([]byte(file), anchor...)
								link.Destination = dest
							}
						}
					}
				}

				return ast.WalkSkipChildren, nil
			}
		}

		return ast.WalkContinue, nil
	})

	md.Renderer().Render(&fixed, source, node)
	return fixed.Bytes(), nil
}

// getParent gets the correct parent of file between topDir and workDir
func getParent(file, topDir, workDir string) (string, error) {
	p := path.Join(workDir, file)
	_, err := os.Stat(p)
	if err == nil {
		return workDir, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	} else if workDir == topDir {
		return "", err
	}

	workDir = path.Clean(path.Join(workDir, ".."))
	return getParent(file, topDir, workDir)
}
