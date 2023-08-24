package gourd

import (
	"fmt"
	"os"
	"time"
)

type (
	StattedBucketer struct {
		StepName string
		Bucketer Bucketer
	}
)

var (
	_ Bucketer = StattedBucketer{}
)

func (b StattedBucketer) Bucket(in Buckets) (Buckets, error) {
	onlyPossiblesBefore := in.PossibleDuplicates()
	var numFilesBefore int
	for _, bucket := range onlyPossiblesBefore {
		numFilesBefore += len(bucket)
	}
	bucketsBefore := len(onlyPossiblesBefore)

	start := time.Now()
	out, err := b.Bucketer.Bucket(onlyPossiblesBefore)
	if err != nil {
		return nil, err
	}
	duration := time.Since(start)
	// TODO: clearing helpful or detrimental?
	clear(onlyPossiblesBefore)

	onlyPossiblesAfter := out.PossibleDuplicates()
	bucketsAfter := len(onlyPossiblesAfter)
	clear(out)

	var numFilesAfter int
	var totalSize int64
	for _, bucket := range onlyPossiblesAfter {
		numFilesAfter += len(bucket)
		totalSize += bucket.TotalFileSize()
	}

	os.Stderr.WriteString(fmt.Sprintf("%11s: Before: %d|%d After: %d|%d Eliminated: %d|%d TotalSize: %s Took: %s\n", b.StepName, bucketsBefore, numFilesBefore, bucketsAfter, numFilesAfter, bucketsBefore-bucketsAfter, numFilesBefore-numFilesAfter, HumanReadableSize(totalSize), duration.String()))
	return onlyPossiblesAfter, nil
}

func HumanReadableSize(size int64) string {
	fsize := float64(size)
	for _, suffix := range []string{"B", "kiB", "MiB", "GiB", "TiB"} {
		if fsize > 1024 {
			fsize /= 1024
		} else {
			return fmt.Sprintf("%.1f%s", fsize, suffix)
		}
	}
	// bigger than 1024 TiB. Report as PiB, even if greater than 1024 PiB because I had to stop somewhere
	return fmt.Sprintf("%.1fPiB", fsize)
}
