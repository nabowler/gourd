package main

import (
	"crypto"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/nabowler/gourd"
)

func main() {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
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

	rootPath := flagSet.Arg(0)

	os.Stderr.WriteString(fmt.Sprintf("PATH: %s\n", rootPath))

	buckets := gourd.Buckets{}
	firstStep := true
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if firstStep {
			firstStep = false
		} else if d.IsDir() && !*recursive {
			return fs.SkipDir
		}

		if err != nil {
			return err
		}
		// os.Stderr.WriteString(fmt.Sprintf("Walking %s: %+v\n", path, d))

		// TODO: configurable/extensible filtering
		if d.IsDir() {
			if d.Name() == ".git" {
				os.Stderr.WriteString(fmt.Sprintf("Skipping Directory %s\n", path))
				return fs.SkipDir
			}

			return nil
		}

		buckets["base"] = append(buckets["base"], path)
		return nil
	})

	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}
	// os.Stderr.WriteString(fmt.Sprintf("Found buckets: %+v\n", buckets))

	chained := gourd.ChainedBucketer{
		Bucketers: []gourd.Bucketer{
			maybeVerbose("Size", gourd.NewFileSizeBucketer(*minFileSize), *verbose),
			// SameFilterBucketer is quadratic per-bucket, so it should be faster to run it
			// after 1-or-more Bucketers to reduce the number of comparisons.
			// TODO: benchmark this, and the effects of ordering
			maybeVerbose("Same File", gourd.SameFilterBucketer{}, *verbose),
			maybeVerbose("First Byte", gourd.NewFirstByteBucketer(), *verbose),
			maybeVerbose("Last Byte", gourd.NewLastByteBucketer(), *verbose),
		},
	}
	if *md5 {
		chained.Bucketers = append(chained.Bucketers, maybeVerbose("MD5", must(func() (gourd.Bucketer, error) {
			return gourd.NewCryptoHashBucketer(crypto.MD5)
		}), *verbose))
	}
	if *sha1 {
		chained.Bucketers = append(chained.Bucketers, maybeVerbose("SHA1", must(func() (gourd.Bucketer, error) {
			return gourd.NewCryptoHashBucketer(crypto.SHA1)
		}), *verbose))
	}
	if *sha256 {
		chained.Bucketers = append(chained.Bucketers, maybeVerbose("SHA256", must(func() (gourd.Bucketer, error) {
			return gourd.NewCryptoHashBucketer(crypto.SHA256)
		}), *verbose))
	}
	if *sha512 {
		chained.Bucketers = append(chained.Bucketers, maybeVerbose("SHA512", must(func() (gourd.Bucketer, error) {
			return gourd.NewCryptoHashBucketer(crypto.SHA512)
		}), *verbose))
	}

	buckets, err = chained.Bucket(buckets)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}

	// os.Stderr.WriteString(fmt.Sprintf("Final Buckets:\n%+v\n", buckets))
	probableDuplicates := buckets.PossibleDuplicates()
	fmt.Printf("Probable Duplicates:\n%s\n", buckets.PossibleDuplicates())

	if *makeHardLinks {
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

func tempFileSuffix() string {
	return fmt.Sprintf(".gourd%d_%d", os.Getpid(), time.Now().UnixNano())
}
