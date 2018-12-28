package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wonsikin/codesearch/cmd/isearch/src"
)

var usageMessage = `usage: cindex [-list] [-reset] [path...]
Cindex prepares the trigram index for use by csearch.  The index is the
file named by $CSEARCHINDEX, or else $HOME/.csearchindex.
The simplest invocation is
	cindex path...
which adds the file or directory tree named by each path to the index.
For example:
	cindex $HOME/src /usr/include
or, equivalently:
	cindex $HOME/src
	cindex /usr/include
If cindex is invoked with no paths, it reindexes the paths that have
already been added, in case the files have changed.  Thus, 'cindex' by
itself is a useful command to run in a nightly cron job.
The -list flag causes cindex to list the paths it has indexed and exit.
By default cindex adds the named paths to the index but preserves 
information about other paths that might already be indexed
(the ones printed by cindex -list).  The -reset flag causes cindex to
delete the existing index before indexing the new paths.
With no path arguments, cindex -reset removes the index.
`

func usage() {
	fmt.Fprintf(os.Stderr, usageMessage)
	os.Exit(2)
}

var (
	fFlag   = flag.String("f", "", "search only files with names matching this regexp")
	iFlag   = flag.Bool("i", false, "case-insensitive search")
	csvFlag = flag.String("csv", "./i18n.csv", "I18N csv")

	matches bool
)

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	fmt.Printf("args is %#+v\n", args)
	fmt.Printf("fFlag is %s\n", *fFlag)

	// create index
	paths := []string{*fFlag}
	src.CreateIndex(paths)

	// read I18N.csv
	content, err := src.Parser(*csvFlag)
	if err != nil {
		log.Fatal("")
		os.Exit(1)
		return
	}

	// search
	src.Search(*fFlag, content)
	// TODO: generate analysis results
}
