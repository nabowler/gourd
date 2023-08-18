# Gourd

A duplicate file finder inspired by rdfind

## Notes:

Creates a chain of `Bucketers` that take a `Bucket` (map[any][]Path) and return a new `Bucket`.

`Bucket`s of a single path can be dropped/ignored.

Final steps:
- List final `Bucket`s with multiple paths
- Create hardlinks if -makehardlinks is set

## Planned Bucketers

- md5
- sha1
- sha256
- sha512
- firstbyte
- lastbyte
- filesize
- mimetype

## Other steps

- device and inode

Used to find existing links ?

symlinks vs hardlinks ?

## References

Make Hardlinks

- https://pkg.go.dev/golang.org/x/sys/unix#Link
- https://man7.org/linux/man-pages/man2/link.2.html

Device and inode

- https://pkg.go.dev/golang.org/x/sys/unix#Fstat
- https://pkg.go.dev/golang.org/x/sys/unix#Stat_t
- https://man7.org/linux/man-pages/man3/fstatat.3p.html