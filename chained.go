package gourd

type (
	// ChainedBucketer applies a series of Bucketers.
	ChainedBucketer struct {
		Bucketers []Bucketer
	}
)

var (
	_ Bucketer = ChainedBucketer{}
)

func (bm ChainedBucketer) Bucket(in Buckets) (Buckets, error) {
	if len(bm.Bucketers) == 0 {
		return nil, nil
	}

	buckets := in
	var err error
	for i := range bm.Bucketers {
		// TODO: clear previous in?
		buckets, err = bm.Bucketers[i].Bucket(buckets)
		if err != nil {
			return nil, err
		}
		// os.Stderr.WriteString(fmt.Sprintf("CurrentBuckets:\n%s", buckets))
	}

	return buckets, nil
}
