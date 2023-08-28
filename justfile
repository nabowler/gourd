default: build

build:
  go build -mod=readonly -ldflags="-s -w" cmd/gourd/gourd.go

install: build
  cp gourd ~/go/bin/

install-github:
  go install github.com/nabowler/gourd/cmd/gourd

bench testDir:
  hyperfine -w 5 -N  --export-markdown /tmp/gourd-rdfind-comparison.md 'gourd -r -md5 {{testDir}}' 'rdfind -makeresultsfile false -checksum md5 -dryrun true {{testDir}}' 'gourd -r -sha1 {{testDir}}' 'rdfind -makeresultsfile false -checksum sha1 -dryrun true {{testDir}}' 'gourd -r -sha256 {{testDir}}' 'gourd -r -sha512 {{testDir}}'