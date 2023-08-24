package gourd

import (
	"io/fs"
	"os"

	"golang.org/x/exp/maps"
)

type (
	// SameFilterBucketer filters out files that appear to already be the same file on disk, as per `os.SameFile`.
	SameFilterBucketer struct{}
)

var (
	_ Bucketer = SameFilterBucketer{}
)

func (bm SameFilterBucketer) Bucket(in Buckets) (Buckets, error) {
	out := Buckets{}
	for currentBucketName, bucket := range in {
		if len(bucket) < 2 {
			continue
		}

		filteredFiles := map[Path]fs.FileInfo{}

		for _, toTest := range bucket {
			toTestFi, err := os.Stat(toTest)
			if err != nil {
				return nil, err
			}
			duplicate := false
			for _, filteredFI := range filteredFiles {
				duplicate = os.SameFile(toTestFi, filteredFI)
				if duplicate {
					break
				}
			}
			if duplicate {
				continue
			}
			filteredFiles[toTest] = toTestFi
		}

		out[currentBucketName] = maps.Keys(filteredFiles)
	}

	return out, nil
}
