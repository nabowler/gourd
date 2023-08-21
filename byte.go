package gourd

import (
	"fmt"
	"os"
)

func NewFirstByteBucketer() Bucketer {
	return fileBucketer{
		subBucketFunc: func(f *os.File) (string, bool, error) {
			firstByte := make([]byte, 1)
			n, err := f.Read(firstByte)
			if err != nil {
				return "", false, fmt.Errorf("Unable to read first byte: %w", err)
			}
			if n != len(firstByte) {
				return "", false, fmt.Errorf("No bytes read")
			}
			return fmt.Sprintf("FirstByte:%x", firstByte), true, nil
		},
	}
}

func NewLastByteBucketer() Bucketer {
	return fileBucketer{
		subBucketFunc: func(f *os.File) (string, bool, error) {
			stat, err := f.Stat()
			if err != nil {
				return "", false, fmt.Errorf("Unable to stat the file: %w", err)
			}

			lastByte := make([]byte, 1)
			n, err := f.ReadAt(lastByte, stat.Size()-1)
			if err != nil {
				return "", false, fmt.Errorf("Unable to read last byte: %w", err)
			}
			if n != len(lastByte) {
				return "", false, fmt.Errorf("No bytes read")
			}
			return fmt.Sprintf("LastByte:%x", lastByte), true, nil
		},
	}
}
