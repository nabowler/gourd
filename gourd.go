package gourd

import (
	"fmt"
	"os"
	"strings"
)

type (
	// File is a path and FileInfo.
	File struct {
		Path     string
		FileInfo os.FileInfo
	}

	// Bucket is a list of Files sharing common attributes
	Bucket []File

	// Buckets are a map of a common attributes to a Bucket. The key is determined by the Bucketer, but must include the current Bucket name.
	Buckets map[string]Bucket

	// Bucketer receives Buckets and returns a set of Buckets.
	// Bucketers should ignore a Bucket with `len(in[key]) < 2`.
	// Similarly, Bucketers may omit a Bucket with less than 2 items in the returned Buckets, but are not required to.
	Bucketer interface {
		Bucket(in Buckets) (Buckets, error)
	}
)

// SubBucketName makes consistent bucket namings.
func SubBucketName(currentBucketName, newBucketName string) string {
	return fmt.Sprintf("%s::%s", currentBucketName, newBucketName)
}

func (b Bucket) TotalFileSize() int64 {
	var totalSize int64
	for _, path := range b {
		totalSize += path.FileInfo.Size()
	}
	return totalSize
}

// PossibleDuplicates returns Buckets containing at least two entries.
func (b Buckets) PossibleDuplicates() Buckets {
	out := Buckets{}
	for k, v := range b {
		if len(v) > 1 {
			out[k] = v
		}
	}
	return out
}

func (b Buckets) String() string {
	sb := new(strings.Builder)
	for k, v := range b {
		sb.WriteString(k)
		sb.WriteRune('\n')
		for _, f := range v {
			sb.WriteString("  ")
			sb.WriteString(f.Path)
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}
