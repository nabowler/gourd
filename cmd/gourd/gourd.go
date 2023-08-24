package main

import (
	"crypto"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nabowler/gourd"
)

type (
	cliConfig struct {
		firstByteSize int64
		lastByteSize  int64
		md5           bool
		sha1          bool
		sha256        bool
		sha512        bool
		minFileSize   int64
		makeHardLinks bool
		recursive     bool
		verbose       bool
		rootPath      string
	}
)

func main() {
	config := parseCommandLineArgs()

	// os.Stderr.WriteString(fmt.Sprintf("PATH: %s\n", rootPath))

	buckets, err := gourd.DirWalker{
		Key: "base",
		Exclude: map[string]struct{}{
			".git": {},
		},
		Recursive: config.recursive,
	}.Walk(config.rootPath)

	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}
	// os.Stderr.WriteString(fmt.Sprintf("Found buckets: %+v\n", buckets))

	bucketers := setupBucketers(config)

	buckets, err = bucketers.Bucket(buckets)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}

	// os.Stderr.WriteString(fmt.Sprintf("Final Buckets:\n%+v\n", buckets))
	probableDuplicates := buckets.PossibleDuplicates()
	fmt.Printf("Probable Duplicates:\n%s\n", buckets.PossibleDuplicates())

	if config.makeHardLinks {
		makeHardLinks(probableDuplicates)
	}
}

func parseCommandLineArgs() cliConfig {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	firstByteSize := flagSet.Int64("firstbytessize", 1, "Number of bytes to check at the start of the file. Must be > 0")
	lastByteSize := flagSet.Int64("lastbytessize", 1, "Number of bytes to check at the end of the file. Must be > 0")
	md5 := flagSet.Bool("md5", false, "Apply MD5 bucketing")
	sha1 := flagSet.Bool("sha1", false, "Apply SHA1 bucketing")
	sha256 := flagSet.Bool("sha256", false, "Apply SHA256 bucketing")
	sha512 := flagSet.Bool("sha512", false, "Apply SHA512 bucketing")
	minFileSize := flagSet.Int64("minfilesize", 1, "Minimum file size in bytes")
	makeHardLinks := flagSet.Bool("makehardlinks", false, "Make hard links of probable-duplicates")
	recursive := flagSet.Bool("r", false, "Recursive")
	verbose := flagSet.Bool("v", false, "Verbose")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		flagSet.Usage()
		os.Exit(1)
	}

	if flagSet.NArg() != 1 {
		os.Stderr.WriteString(fmt.Sprintf("Usage: %s PATH\n", os.Args[0]))
		os.Exit(1)
	}

	if !*md5 && !*sha1 && !*sha256 && !*sha512 {
		*sha1 = true
	}

	if *firstByteSize <= 0 {
		os.Stderr.WriteString(fmt.Sprintf("Invalid firstbytessize: %d\n", *firstByteSize))
		flagSet.Usage()
		os.Exit(1)
	}

	if *lastByteSize <= 0 {
		os.Stderr.WriteString(fmt.Sprintf("Invalid lastbytesize: %d\n", *lastByteSize))
		flagSet.Usage()
		os.Exit(1)
	}

	rootPath := flagSet.Arg(0)
	return cliConfig{
		firstByteSize: *firstByteSize,
		lastByteSize:  *lastByteSize,
		md5:           *md5,
		sha1:          *sha1,
		sha256:        *sha256,
		sha512:        *sha512,
		minFileSize:   *minFileSize,
		makeHardLinks: *makeHardLinks,
		recursive:     *recursive,
		verbose:       *verbose,
		rootPath:      rootPath,
	}
}

func setupBucketers(config cliConfig) gourd.Bucketer {
	chained := gourd.ChainedBucketer{
		Bucketers: []gourd.Bucketer{
			maybeVerbose("Size", gourd.NewFileSizeBucketer(config.minFileSize), config.verbose),
			// SameFilterBucketer is quadratic per-bucket, so it should be faster to run it
			// after 1-or-more Bucketers to reduce the number of comparisons.
			// TODO: benchmark this, and the effects of ordering
			// We're not attempting device/inode specific logic because we're using the stdlib os.SameFile
			// functionality which is OS agnostic for us to get a bool, but OS specific in its implementation
			maybeVerbose("Same File", gourd.SameFilterBucketer{}, config.verbose),
			maybeVerbose("First Bytes", gourd.NewFirstBytesBucketer(config.firstByteSize), config.verbose),
			maybeVerbose("Last Bytes", gourd.NewLastBytesBucketer(config.lastByteSize), config.verbose),
		},
	}
	if config.md5 {
		chained.Bucketers = append(chained.Bucketers, maybeVerbose("MD5", must(func() (gourd.Bucketer, error) {
			return gourd.NewCryptoHashBucketer(crypto.MD5)
		}), config.verbose))
	}
	if config.sha1 {
		chained.Bucketers = append(chained.Bucketers, maybeVerbose("SHA1", must(func() (gourd.Bucketer, error) {
			return gourd.NewCryptoHashBucketer(crypto.SHA1)
		}), config.verbose))
	}
	if config.sha256 {
		chained.Bucketers = append(chained.Bucketers, maybeVerbose("SHA256", must(func() (gourd.Bucketer, error) {
			return gourd.NewCryptoHashBucketer(crypto.SHA256)
		}), config.verbose))
	}
	if config.sha512 {
		chained.Bucketers = append(chained.Bucketers, maybeVerbose("SHA512", must(func() (gourd.Bucketer, error) {
			return gourd.NewCryptoHashBucketer(crypto.SHA512)
		}), config.verbose))
	}
	return chained
}

func maybeVerbose(name string, bucketer gourd.Bucketer, verbose bool) gourd.Bucketer {
	if !verbose {
		return bucketer
	}

	return gourd.StattedBucketer{
		StepName: name,
		Bucketer: bucketer,
	}
}

func must[T any](f func() (T, error)) T {
	t, err := f()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}
	return t
}

func makeHardLinks(probableDuplicates gourd.Buckets) {
	suffix := tempFileSuffix()
	for _, bucket := range probableDuplicates {
		masterFilePath := bucket[0]
		for i := 1; i < len(bucket); i++ {
			oldPath := bucket[i]
			tempPath := oldPath + suffix
			err := os.Rename(oldPath, tempPath)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Error renaming duplicate file %s: %v\n", oldPath, err))
				os.Exit(1)
			}
			if err = os.Link(masterFilePath, oldPath); err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Error linking duplicate file %s: %v\n", oldPath, err))
				if err = os.Rename(tempPath, oldPath); err != nil {
					os.Stderr.WriteString(fmt.Sprintf("Error renaming duplicate temp file %s after previous error: %v\n", tempPath, err))
				}
				os.Exit(1)
			}
			if err = os.Remove(tempPath); err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Error removing duplicate temp file %s after hardlink: %v\n", tempPath, err))
				os.Exit(1)
			}
		}
	}
}

func tempFileSuffix() string {
	return fmt.Sprintf(".gourd%d_%d", os.Getpid(), time.Now().UnixNano())
}
