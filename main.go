package main

import (
	"log"
	"os"

	"github.com/henderjm/word-frequency/counter"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	NumberOfWords int    `short:"n" long:"number-of-words" description:"Report a number of most common words on a wiki page; it must be positive" required:"true"`
	PageIDs       string `short:"i" long:"page-ids" description:"page id of a wiki page" required:"true"`
}

func main() {
	var opts Options
	var parser = flags.NewParser(&opts, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			parser.WriteHelp(os.Stderr)
			os.Exit(1)
		}
	}
	if opts.NumberOfWords <= 0 {
		log.Fatalf("n must be greater than 0")
	}
	wf := counter.NewWordFrequencyCounter(opts.NumberOfWords, opts.PageIDs)
	err := wf.Run()
	if err != nil {
		log.Fatalf("could not get data for wiki id %s\n%s", opts.PageIDs, err.Error())
	}
	wf.Write(os.Stdout)
}
