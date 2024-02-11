package main

import (
	"flag"
	"io"
	"os"
	"path"

	"gitlab.com/mcdonagh/mdfix/core"
	"gitlab.com/mcdonagh/mdfix/fixers"
)

func main() {
	in_place := flag.Bool("i", false, "in place")
	textwidth := flag.Int("w", 80, "text width")

	flag.Parse()
	fileName := flag.Arg(0)
	var file *os.File
	var err error
	if fileName == "" {
		file = os.Stdin
	} else {
		file, err = os.Open(fileName)
		if err != nil {
			panic(err)
		}
	}

	source, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	if *in_place {
		file, err = os.Create(fileName)
		if err != nil {
			panic(err)
		}
	} else {
		file = os.Stdout
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	err = core.Fix(source, file, fixers.Options{
		TextWidth: *textwidth,
		TopDir:    cwd,
		WorkDir:   path.Dir(fileName),
	})

	if err != nil {
		panic(err)
	}
}
