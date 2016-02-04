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

// MockableBucket is the interface to both s3.Bucket and
// s3mock.mockBucket, which means it is mockable.  In particular, to
// make a function which accesses S3 service testable, please make
// sure that it doesn't use github.com/AdRoll/goamz/s3.Bucket, but
// github.com/topicai/s3.Bucket.  For example, please refer to
// s3_test.go.
type MockableBucket interface {
	GetReader(p string) (io.ReadCloser, error)
	PutReader(p string, r io.Reader, length int64, contType string, perm s3.ACL, options s3.Options) error
}

type MockS3 struct {
	*s3.S3
}

type bucket struct {
	*s3.Bucket
	Fs map[string][]byte
}

func NewMock(auth aws.Auth, region aws.Region) *MockS3 {
	return &MockS3{
		S3: &s3.S3{Auth: auth, Region: region},
	}
}

func (src *MockS3) Bucket(name string) *bucket {
	return &bucket{
		Bucket: &s3.Bucket{S3: src.S3, Name: name},
		Fs:     make(map[string][]byte),
	}
}

// GetReader retrieves an object from an S3 bucket, returning the body
// of the HTTP response.  It is the caller's responsibility to call
// Close on rc when finished reading.
func (bkt *bucket) GetReader(p string) (io.ReadCloser, error) {
	if f, ok := bkt.Fs[path.Clean(p)]; ok {
		return ioutil.NopCloser(bytes.NewReader(f)), nil
	}
	return nil, fmt.Errorf("File %s not exist in bucket %s", path.Clean(p), bkt.Bucket.Name)
}

func (bkt *bucket) PutReader(p string, r io.Reader, length int64, contType string, perm s3.ACL, options s3.Options) error {
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
