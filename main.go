package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"gitlab.com/mcdonagh/mdfix/core"
	"gitlab.com/mcdonagh/mdfix/fixers"
)

func main() {
	in_place := flag.Bool("i", false, "in place")
	text_width := flag.Int("w", 80, "text width")
	n_parallel := flag.Int("n", 16, "number of parallel fixes")

	flag.Parse()
	file_chan := make(chan string)
	var wg sync.WaitGroup
	for i := 0; i < *n_parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file_name := range file_chan {
				file, err := os.Open(file_name)
				if err != nil {
					panic(err)
				}

				source, err := io.ReadAll(file)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}

				err = file.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}

				if *in_place {
					file, err = os.Create(file_name)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return
					}
				} else {
					file = os.Stdout
				}

				cwd, err := os.Getwd()
				if err != nil {
					cwd = ""
				}

				err = core.Fix(source, file, fixers.Options{
					TextWidth: *text_width,
					TopDir:    cwd,
					WorkDir:   path.Dir(file_name),
				})

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}
			}
		}()
	}

	for _, file_name := range flag.Args() {
		file_chan <- file_name
	}

	close(file_chan)
	wg.Wait()
}
