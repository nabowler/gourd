# Gourd

Gourd is a command line tool to find duplicate files.

Gourd is inspired by [rdfind](https://github.com/pauldreik/rdfind), but is not designed to be compatible with rdfind in terms of output or command line flags.

Gourd came from a usecase where I wanted to deduplicate data on a server, but could not install `rdfind` natively and could not use `rdfind` from my local machine due to differing libc versions. I was able to work around this problem using Docker, but it is unwieldy and cumbersome to do so.

## Cautions and Warnings

This software is experimental and untested.

**Use at your own risk.**

## Build

```sh
git clone https://github.com/nabowler/gourd.git
go build -mod=readonly -ldflags="-s -w" cmd/gourd/gourd.go
```

## Installation

With Go
```bash
go install github.com/nabowler/gourd/cmd/gourd
```

## Use

```bash
gourd -r -v -sha1 path/to/directory
```

## Benchmarks and Comparison to rdfind

TODO

## Development Notes
### Planned Bucketers

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

### Other steps

- [x] duplicate device and inode detection

### To Consider

- goroutines for Hash/file steps?
  - one per in bucket?
- -exclude patterns
  - for instance, could useful when only wanting to deduplicate .jpg files in a mixed-media directory tree
- multiple paths/files as args (currently only accepts 1)