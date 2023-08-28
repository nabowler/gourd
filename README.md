# Gourd

Gourd is a command line tool to find duplicate files.

## Acknowledgements

Gourd is inspired by my use of [rdfind](https://github.com/pauldreik/rdfind), but is not designed to be compatible with rdfind in terms of output or command line flags. Gourd is not related to, a port of, or based on the source of `rdfind`.

Gourd came from a usecase where I wanted to deduplicate data on a server, but could not install `rdfind` natively and could not use `rdfind` from my local machine due to differing libc versions. I was able to work around this problem using Docker, but it is unwieldy and cumbersome to do so, and wanted an easily-portable solution.

# Cautions and Warnings

This software is experimental and untested.

**Use at your own risk.**

# Build

```sh
git clone https://github.com/nabowler/gourd.git
go build -mod=readonly -ldflags="-s -w" cmd/gourd/gourd.go
```

# Installation

With Go
```bash
go install github.com/nabowler/gourd/cmd/gourd
```

# Use

```bash
gourd -r -v -sha1 path/to/directory
```

# Benchmarks and Comparison to rdfind

## Benchmarks

### 11.7 GiB images containing 184 duplicate files

```sh
$ hyperfine -w 5 -N  --export-markdown ~/gourd-rdfind-comparison.md 'gourd -r -md5 .' 'rdfind -makeresultsfile false -checksum md5 -dryrun true .'
Benchmark 1: gourd -r -md5 .
  Time (mean ± σ):     893.8 ms ±  48.4 ms    [User: 720.4 ms, System: 198.0 ms]
  Range (min … max):   840.0 ms … 1005.1 ms    10 runs

Benchmark 2: rdfind -makeresultsfile false -checksum md5 -dryrun true .
  Time (mean ± σ):      1.025 s ±  0.018 s    [User: 0.866 s, System: 0.155 s]
  Range (min … max):    1.001 s …  1.047 s    10 runs

Summary
  gourd -r -md5 . ran
    1.15 ± 0.07 times faster than rdfind -makeresultsfile false -checksum md5 -dryrun true .
```

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `gourd -r -md5 .` | 893.8 ± 48.4 | 840.0 | 1005.1 | 1.00 |
| `rdfind -makeresultsfile false -checksum md5 -dryrun true .` | 1025.2 ± 18.0 | 1000.6 | 1047.1 | 1.15 ± 0.07 |


Comparison of findings

```sh
$ gourd -r -v -md5 . > /dev/null
Found 9286 files totaling 11.7GiB
       Size: Before: 1|9286 After: 115|232 Eliminated: -114|9054 TotalSize: 447.6MiB Took: 139.169211ms
  Same File: Before: 115|232 After: 115|232 Eliminated: 0|0 TotalSize: 447.6MiB Took: 55.478µs
First Bytes: Before: 115|232 After: 94|190 Eliminated: 21|42 TotalSize: 435.2MiB Took: 2.647062ms
 Last Bytes: Before: 94|190 After: 91|184 Eliminated: 3|6 TotalSize: 434.1MiB Took: 2.074249ms
        MD5: Before: 91|184 After: 91|184 Eliminated: 0|0 TotalSize: 434.1MiB Took: 638.124048ms
Found 184 duplicate files with 220.1MiB reclaimable space
```

```sh
$ rdfind -makeresultsfile false -checksum md5 -dryrun true .                                                                                                       (DRYRUN MODE) Now scanning ".", found 9286 files.
(DRYRUN MODE) Now have 9286 files in total.
(DRYRUN MODE) Removed 0 files due to nonunique device and inode.
(DRYRUN MODE) Total size is 12509236983 bytes or 12 GiB
Removed 9054 files due to unique sizes from list. 232 files left.
(DRYRUN MODE) Now eliminating candidates based on first bytes: removed 42 files from list. 190 files left.
(DRYRUN MODE) Now eliminating candidates based on last bytes: removed 6 files from list. 184 files left.
(DRYRUN MODE) Now eliminating candidates based on md5 checksum: removed 0 files from list. 184 files left.
(DRYRUN MODE) It seems like you have 184 files that are not unique
(DRYRUN MODE) Totally, 220 MiB can be reduced.
```

## Comparisons

```sh
$ ldd $(which rdfind)
        linux-vdso.so.1 (0x00007ffd8ba03000)
        libnettle.so.8 => /usr/lib/libnettle.so.8 (0x00007fec79fb7000)
        libstdc++.so.6 => /usr/lib/libstdc++.so.6 (0x00007fec79c00000)
        libgcc_s.so.1 => /usr/lib/libgcc_s.so.1 (0x00007fec79f8f000)
        libc.so.6 => /usr/lib/libc.so.6 (0x00007fec79800000)
        libm.so.6 => /usr/lib/libm.so.6 (0x00007fec79e9f000)
        /lib64/ld-linux-x86-64.so.2 => /usr/lib64/ld-linux-x86-64.so.2 (0x00007fec7a06f000)

$ ldd $(which gourd)
        not a dynamic executable
```

```sh
$ du -h $(which rdfind)
96K     /usr/bin/rdfind

$ du -h $(which gourd)
2.3M    /home/nathan/go/bin/gourd
```



# Development Notes
## Planned Bucketers

- [x] md5
  - configurable from commandline (-md5 flag)
- [x] sha1
  - configurable from commandline (-sha1 flag)
- [x] sha256
  - configurable from commandline (-sha256 flag)
- [x] sha512
  - configurable from commandline (-sha512 flag)
- [x] firstbytes
  - number of bytes configurable with -firstbytessize flag
- [x] lastbytes
  - number of bytes configurable with -lastbytessize flag
- [x] filesize
- [x] statted
  - Outputs information about the number of files before and after an inner Bucketer
  - configurable from commandline (-v flag)

Note: -md5, -sha1, -sha256, and -sha512 are additive and will be applied in that order if set. If none are set, a default of SHA-1 is used.

## Other steps

- [x] duplicate device and inode detection

## To Consider

- goroutines for Hash/file steps?
  - one per in bucket?
- -exclude patterns
  - for instance, could useful when only wanting to dedupwhilicate .jpg files in a mixed-media directory tree
