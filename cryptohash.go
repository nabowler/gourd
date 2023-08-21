package gourd

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"fmt"
	"io"
	"os"
)

// NewCryptoHashBucketer returns a Bucketer that calcuates the provided crypto.Hash on the file.
// By default, MD5, SHA1, SHA256, and SHA512 are supported. Other hashes are only supported if `hash.Available()` is true.
func NewCryptoHashBucketer(hash crypto.Hash) (Bucketer, error) {
	fbm := fileBucketer{
		subBucketFunc: func(f *os.File) (string, bool, error) {
			hasher := hash.New()
			_, err := io.Copy(hasher, f)
			if err != nil {
				return "", false, fmt.Errorf("Unable to hash file %w", err)
			}

			return fmt.Sprintf("%s:%x", hash.String(), hasher.Sum(nil)), true, nil
		},
	}

	if !hash.Available() {
		return fbm, fmt.Errorf("Hash %s is not available", hash.String())
	}
	return fbm, nil
}
