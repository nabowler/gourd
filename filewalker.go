package gourd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"syscall"
)

type (
	DirWalker struct {
		Key         string
		Exclude     map[string]struct{}
		Recursive   bool
		AppendDevID bool
	}
)

func (dw DirWalker) Walk(rootPaths ...string) (Buckets, error) {
	buckets := Buckets{}
	exploredPaths := map[string]struct{}{}

	for _, rootPath := range rootPaths {
		key := dw.Key
		if key == "" {
			key = rootPath
		}
		firstStep := true
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

			if _, ok := exploredPaths[path]; ok {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}

			exploredPaths[path] = struct{}{}

			if d.IsDir() {
				return nil
			}

			fileInfo, err := d.Info()
			if err != nil {
				return err
			}
			key := key
			if dw.AppendDevID {
				sys := fileInfo.Sys()
				switch t := sys.(type) {
				case *syscall.Stat_t:
					key = fmt.Sprintf("%s::%d", key, t.Dev)
				default:
					fmt.Println(reflect.TypeOf(sys))
				}
			}

			buckets[key] = append(buckets[key], File{
				Path:     path,
				FileInfo: fileInfo,
			})
			return nil
		})
		if err != nil {
			return buckets, err
		}
	}

	return buckets, nil
}
