package gourd

import (
	"os"
)

type (
	// SameFilterBucketer filters out files that appear to already be the same file on disk, as per `os.SameFile`.
	// Unlike most other Bucketer implementations, this will only remove entries from Buckets rather than create
	// new Buckets.
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

		filteredFiles := []File{}

		for _, toTest := range bucket {
			toTestFi := toTest.FileInfo
			duplicate := false
			for _, filteredFI := range filteredFiles {
				duplicate = os.SameFile(toTestFi, filteredFI.FileInfo)
				if duplicate {
					break
				}
			}
			if duplicate {
				continue
			}
			filteredFiles = append(filteredFiles, toTest)
		}

		out[currentBucketName] = filteredFiles
	}

	return out, nil
}
