package gourd

import (
	"io/fs"
	"path/filepath"
)

type (
	DirWalker struct {
		Key       string
		Exclude   map[string]struct{}
		Recursive bool
	}
)

func (dw DirWalker) Walk(rootPath string) (Buckets, error) {
	firstStep := true
	buckets := Buckets{}
	key := dw.Key
	if key == "" {
		key = rootPath
	}
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if firstStep {
			// on the first step in the walk, we do need to "recurse" into the directory
			firstStep = false
		} else if d.IsDir() && !dw.Recursive {
			// all other steps, only recurse when resursive is true
			return fs.SkipDir
		}

		if err != nil {
			return err
		}

		if _, ok := dw.Exclude[path]; ok {
			if d.IsDir() {
				// if exclusion matches a directory, skip the entire directory
				return fs.SkipDir
			}
			// exclusion matches a file, just ignore it and continue on
			return nil
		}

		if d.IsDir() {
			return nil
		}

		buckets[key] = append(buckets[key], path)
		return nil
	})

	return buckets, err
}
