package gourd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

// NewFirstBytesBucketer creates a Bucketer that examines the first `numBytes` of a file.
// If numBytes is more than the file size, the file size is used as numBytes.
// If numBytes is <= 0 after comparing to the file size, it is bucketed under a constant key
// that will not overlap with normal bucket keys
func NewFirstBytesBucketer(numBytes int64) Bucketer {
	return fileBucketer{
		subBucketFunc: func(f *os.File) (string, bool, error) {
			stat, err := f.Stat()
			if err != nil {
				return "", false, fmt.Errorf("Unable to stat the file: %w", err)
			}
			numBytes = min(numBytes, stat.Size())
			if numBytes <= 0 {
				return "FirstBytes:-", true, nil
			}

			firstBytes := make([]byte, numBytes)
			n, err := f.Read(firstBytes)
			if err != nil && !errors.Is(err, io.EOF) {
				return "", false, fmt.Errorf("Unable to read first byte: %w", err)
			}
			if n != len(firstBytes) {
				return "", false, fmt.Errorf("No bytes read")
			}
			return fmt.Sprintf("FirstBytes:%s", base64.StdEncoding.EncodeToString(firstBytes)), true, nil
		},
	}
}

// NewLastBytesBucketer creates a Bucketer that examines the last `numBytes` of a file.
// If numBytes is more than the file size, the file size is used as numBytes.
// If numBytes is <= 0 after comparing to the file size, it is bucketed under a constant key
// that will not overlap with normal bucket keys
func NewLastBytesBucketer(numBytes int64) Bucketer {
	return fileBucketer{
		subBucketFunc: func(f *os.File) (string, bool, error) {
			stat, err := f.Stat()
			if err != nil {
				return "", false, fmt.Errorf("Unable to stat the file: %w", err)
			}

			numBytes = min(numBytes, stat.Size())
			if numBytes <= 0 {
				return "LastBytes:-", true, nil
			}

			lastBytes := make([]byte, numBytes)
			n, err := f.ReadAt(lastBytes, stat.Size()-numBytes)
			if err != nil && !errors.Is(err, io.EOF) {
				return "", false, fmt.Errorf("Unable to read last byte: %w", err)
			}
			if n != len(lastBytes) {
				return "", false, fmt.Errorf("No bytes read")
			}
			return fmt.Sprintf("LastBytes:%s", base64.StdEncoding.EncodeToString(lastBytes)), true, nil
		},
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
