package main

import (
	"flag"
	"io"
	"os"

	"gitlab.com/mcdonagh/mdfix/core"
	"gitlab.com/mcdonagh/mdfix/fixers"
)

func main() {
	in_place := flag.Bool("i", false, "inplace")
	flag.Parse()
	path := flag.Arg(0)
	var file *os.File
	var err error
	if path == "" {
		file = os.Stdin
	} else {
		file, err = os.Open(path)
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
		file, err = os.Create(path)
		if err != nil {
			panic(err)
		}
	} else {
		file = os.Stdout
	}

	err = core.Fix(source, file, fixers.Options{
		TextWidth: 80,
	})

	if err != nil {
		panic(err)
	}
}
