package src

import (
	"encoding/json"
	"log"
	"os"
	"sort"

	"github.com/wonsikin/codesearch/index"
	"github.com/wonsikin/codesearch/regexp"
)

// Record represents a query record
type Record struct {
	Query   string       `json:"query,omitempty"`
	Count   int          `json:"count,omitempty"`
	Grep    *regexp.Grep `json:"-"`
	Content *IWriter     `json:"content,omitempty"`
}

var prefix string

// Search searches in content
func Search(filePrefix string, content [][]string) error {
	prefix = filePrefix
	size := len(content)

	jobs := make(chan string, size)
	results := make(chan *Record, size)

	for w := 1; w <= 10; w++ {
		go worker(w, jobs, results)
	}

	for _, row := range content {
		if len(row) >= 1 {
			jobs <- row[0]
		}
	}
	close(jobs)

	records := make([]*Record, 0)
	for a := 1; a <= size; a++ {
		rc := <-results
		records = append(records, rc)
	}

	sort.Sort(Records(records))
	for _, record := range records {
		log.Printf("query: %s, %d hit\n", record.Query, record.Count)
	}

	of, err := os.Create("./result.json")
	if err != nil {
		log.Fatalf("error caught when creating file: %s", err)
		return err
	}
	defer of.Close()

	data, err := json.Marshal(records)
	if err != nil {
		log.Fatalf("error caught when marshalling data: %s", err)
		return err
	}

	_, err = of.Write(data)
	if err != nil {
		log.Fatalf("error caught when writing data: %s", err)
		return err
	}
	return nil
}

func worker(id int, jobs <-chan string, results chan<- *Record) {
	for j := range jobs {
		log.Printf("worker: %d start to search %s", id, j)
		rec, err := csearch(prefix, j)
		if err != nil {
			log.Fatalf("error caught when executing csearch: %s", err)
			continue
		}

		results <- rec
	}
}

func csearch(filePrefix, query string) (*Record, error) {
	writer := &IWriter{}
	g := regexp.Grep{
		Stdout: writer,
		Stderr: os.Stderr,
	}
	// g.AddFlags()

	pat := "(?m)" + query
	// if *iFlag {
	// 	pat = "(?i)" + pat
	// }
	re, err := regexp.Compile(pat)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	g.Regexp = re
	var fre *regexp.Regexp
	if filePrefix != "" {
		fre, err = regexp.Compile(filePrefix)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	q := index.RegexpQuery(re.Syntax)
	// if *verboseFlag {
	// 	log.Printf("query: %s\n", q)
	// }

	ix := index.Open(index.File())
	var post []uint32
	post = ix.PostingQuery(q)

	if fre != nil {
		fnames := make([]uint32, 0, len(post))

		for _, fileid := range post {
			name := ix.Name(fileid)
			if fre.MatchString(name, true, true) < 0 {
				continue
			}
			fnames = append(fnames, fileid)
		}

		post = fnames
	}

	for _, fileid := range post {
		name := ix.Name(fileid)
		g.File(name)
	}

	record := &Record{
		Query:   query,
		Count:   len(post),
		Grep:    &g,
		Content: writer,
	}

	return record, nil
}
