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
go install github.com/nabowler/gourd/cmd/gourd@latest
```

# Use

```bash
gourd -r -v -sha1 path/to/directory [path/to/directory2 ...]
```

# Benchmarks and Comparison to rdfind

## Benchmarks

### 11.7 GiB images containing 184 duplicate files

```sh
$ hyperfine -w 5 -N  --export-markdown ~/gourd-rdfind-comparison.md 'gourd -r -md5 .' 'rdfind -makeresultsfile false -checksum md5 -dryrun true .' 'gourd -r -sha1 .' 'rdfind -makeresultsfile false -checksum sha1 -dryrun true .' 'gourd -r -sha256 .' 'gourd -r -sha512 .'
Benchmark 1: gourd -r -md5 .
  Time (mean ± σ):     920.6 ms ±  45.0 ms    [User: 728.4 ms, System: 216.5 ms]
  Range (min … max):   835.5 ms … 989.4 ms    10 runs

Benchmark 2: rdfind -makeresultsfile false -checksum md5 -dryrun true .
  Time (mean ± σ):      1.028 s ±  0.016 s    [User: 0.872 s, System: 0.152 s]
  Range (min … max):    1.007 s …  1.049 s    10 runs

Benchmark 3: gourd -r -sha1 .
  Time (mean ± σ):      1.301 s ±  0.058 s    [User: 1.126 s, System: 0.200 s]
  Range (min … max):    1.193 s …  1.383 s    10 runs

Benchmark 4: rdfind -makeresultsfile false -checksum sha1 -dryrun true .
  Time (mean ± σ):      1.246 s ±  0.014 s    [User: 1.084 s, System: 0.158 s]
  Range (min … max):    1.224 s …  1.264 s    10 runs

Benchmark 5: gourd -r -sha256 .
  Time (mean ± σ):      2.636 s ±  0.053 s    [User: 2.446 s, System: 0.216 s]
  Range (min … max):    2.555 s …  2.711 s    10 runs

Benchmark 6: gourd -r -sha512 .
  Time (mean ± σ):      1.915 s ±  0.047 s    [User: 1.717 s, System: 0.222 s]
  Range (min … max):    1.830 s …  1.987 s    10 runs

Summary
  gourd -r -md5 . ran
    1.12 ± 0.06 times faster than rdfind -makeresultsfile false -checksum md5 -dryrun true .
    1.35 ± 0.07 times faster than rdfind -makeresultsfile false -checksum sha1 -dryrun true .
    1.41 ± 0.09 times faster than gourd -r -sha1 .
    2.08 ± 0.11 times faster than gourd -r -sha512 .
    2.86 ± 0.15 times faster than gourd -r -sha256 .
```

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `gourd -r -md5 .` | 920.6 ± 45.0 | 835.5 | 989.4 | 1.00 |
| `rdfind -makeresultsfile false -checksum md5 -dryrun true .` | 1028.2 ± 15.6 | 1007.4 | 1049.3 | 1.12 ± 0.06 |
| `gourd -r -sha1 .` | 1300.9 ± 57.9 | 1193.3 | 1383.2 | 1.41 ± 0.09 |
| `rdfind -makeresultsfile false -checksum sha1 -dryrun true .` | 1246.0 ± 13.9 | 1224.3 | 1264.4 | 1.35 ± 0.07 |
| `gourd -r -sha256 .` | 2636.1 ± 52.9 | 2555.3 | 2711.3 | 2.86 ± 0.15 |
| `gourd -r -sha512 .` | 1914.6 ± 46.6 | 1830.3 | 1987.0 | 2.08 ± 0.11 |

## Comparisons

### Duplicates Found
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

### Linked Dependencies
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

### Binary Size

```sh
$ du -h $(which rdfind)
96K     /usr/bin/rdfind

$ du -h $(which gourd)
2.3M    /home/nathan/go/bin/gourd
```

Note: rdfind installed from system repos. gourd install via `go install`



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
