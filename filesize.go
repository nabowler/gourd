package gourd

import (
	"fmt"
	"os"
)

// NewFileSizeBucketer returns a Bucketer based on the size of the File.
func NewFileSizeBucketer(minSize int64) Bucketer {
	return fileBucketer{
		subBucketFunc: func(f *os.File) (string, bool, error) {
			stat, err := f.Stat()
			if err != nil {
				return "", false, fmt.Errorf("Unable to stat the file: %w", err)
			}

			size := stat.Size()
			if size < minSize {
				return "", false, nil
			}

			return fmt.Sprintf("Size:%d", size), true, nil
		},
	}
}
