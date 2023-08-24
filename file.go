package gourd

import (
	"fmt"
	"os"
)

type (
	// SubBucketFunc determines the sub-bucket name and if the file should be sub-bucketed
	// based on the contents of the file.
	// SubBucketFunc implementations _should not_ Close the *os.File.
	SubBucketFunc func(*os.File) (subBucketName string, accept bool, err error)

	fileBucketer struct {
		subBucketFunc SubBucketFunc
	}
)

var (
	_ Bucketer = fileBucketer{}
)

// NewFileBucketer returns a Bucketer that uses the provided SubBucketFunc to generate the output Buckets.
func NewFileBucketer(sbf SubBucketFunc) (Bucketer, error) {
	fbm := fileBucketer{subBucketFunc: sbf}
	if sbf == nil {
		return fbm, fmt.Errorf("SubBucketFunc is nil")
	}
	return fbm, nil
}

func (bm fileBucketer) Bucket(in Buckets) (Buckets, error) {
	sbf := bm.subBucketFunc
	if sbf == nil {
		// something's gone horribly wrong since none of this is exported
		return nil, fmt.Errorf("SubBucketFunc is nil")
	}

	out := Buckets{}
	for currentBucketName, bucket := range in {
		if len(bucket) < 2 {
			continue
		}

		for i := range bucket {
			path := bucket[i]
			f, err := os.Open(path.Path)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Unable to open %s: %v\n", path, err))
				return nil, err
			}

			subBucketName, accept, err := sbf(f)
			if err != nil {
				_ = f.Close()
				os.Stderr.WriteString(fmt.Sprintf("Unable to process %s: %v\n", path, err))
				return nil, err
			}

			if !accept {
				err = f.Close()
				if err != nil {
					os.Stderr.WriteString(fmt.Sprintf("Unable to close %s: %v\n", path, err))
					return nil, err
				}
				continue
			}

			newBucketName := SubBucketName(currentBucketName, subBucketName)
			out[newBucketName] = append(out[newBucketName], path)

			err = f.Close()
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Unable to close %s: %v\n", path, err))
				return nil, err
			}
		}
	}

	return out, nil
}
