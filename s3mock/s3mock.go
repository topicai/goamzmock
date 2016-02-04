package s3mock

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"

	"github.com/AdRoll/goamz/aws"
	"github.com/AdRoll/goamz/s3"
)

type S3 struct {
	*s3.S3
}

type Bucket struct {
	*s3.Bucket
	Fs map[string][]byte
}

func New(auth aws.Auth, region aws.Region) *S3 {
	return &S3{
		S3: &s3.S3{Auth: auth, Region: region},
	}
}

func (src *S3) Bucket(name string) *Bucket {
	return &Bucket{
		Bucket: &s3.Bucket{S3: src.S3, Name: name},
		Fs:     make(map[string][]byte),
	}
}

// GetReader retrieves an object from an S3 bucket, returning the body
// of the HTTP response.  It is the caller's responsibility to call
// Close on rc when finished reading.
func (bkt *Bucket) GetReader(p string) (io.ReadCloser, error) {
	if f, ok := bkt.Fs[path.Clean(p)]; ok {
		return ioutil.NopCloser(bytes.NewReader(f)), nil
	}
	return nil, fmt.Errorf("File %s not exist in bucket %s", path.Clean(p), bkt.Bucket.Name)
}

func (bkt *Bucket) PutReader(p string, r io.Reader, length int64, contType string, perm s3.ACL, options s3.Options) error {
	b, e := ioutil.ReadAll(r)
	if e != nil {
		return e
	}
	if int64(len(b)) != length {
		return fmt.Errorf("len(b) (%d) != length (%d)", len(b), length)
	}

	bkt.Fs[path.Clean(p)] = b
	return nil
}
