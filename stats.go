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
	sizeBefore := len(onlyPossiblesBefore)

	start := time.Now()
	out, err := b.Bucketer.Bucket(onlyPossiblesBefore)
	if err != nil {
		return nil, err
	}
	duration := time.Since(start)
	// TODO: clearing helpful or detrimental?
	clear(onlyPossiblesBefore)

	onlyPossiblesAfter := out.PossibleDuplicates()
	sizeAfter := len(onlyPossiblesAfter)
	clear(out)

	var numFilesAfter int
	for _, bucket := range onlyPossiblesAfter {
		numFilesAfter += len(bucket)
	}

	os.Stderr.WriteString(fmt.Sprintf("%11s: Before: %d|%d After: %d|%d Eliminated: %d|%d Took: %s\n", b.StepName, sizeBefore, numFilesBefore, sizeAfter, numFilesAfter, sizeBefore-sizeAfter, numFilesBefore-numFilesAfter, duration.String()))
	return onlyPossiblesAfter, nil
}
