package src

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/wonsikin/codesearch/index"
)

// CreateIndex creates index
func CreateIndex(paths []string) {
	if len(paths) == 0 {
		ix := index.Open(index.File())
		for _, arg := range ix.Paths() {
			paths = append(paths, arg)
		}
	}

	// Translate paths to absolute paths so that we can
	// generate the file list in sorted order.
	for i, arg := range paths {
		a, err := filepath.Abs(arg)
		if err != nil {
			log.Printf("%s: %s", arg, err)
			paths[i] = ""
			continue
		}
		paths[i] = a
	}
	sort.Strings(paths)

	for len(paths) > 0 && paths[0] == "" {
		paths = paths[1:]
	}

	master := index.File()
	file := master
	file += "~"

	ix := index.Create(file)
	ix.AddPaths(paths)
	for _, arg := range paths {
		log.Printf("index %s", arg)
		filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if _, elem := filepath.Split(path); elem != "" {
				// Skip various temporary or "hidden" files or directories.
				if elem[0] == '.' || elem[0] == '#' || elem[0] == '~' || elem[len(elem)-1] == '~' {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
			if err != nil {
				log.Printf("%s: %s", path, err)
				return nil
			}
			if info != nil && info.Mode()&os.ModeType == 0 {
				ix.AddFile(path)
			}
			return nil
		})
	}
	log.Printf("flush index")
	ix.Flush()

	log.Printf("merge %s %s", master, file)
	index.Merge(file+"~", master, file)
	os.Remove(file)
	os.Rename(file+"~", master)
	log.Printf("done")
	return
}
