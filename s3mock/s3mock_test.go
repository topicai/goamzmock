package s3mock

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/AdRoll/goamz/aws"
	"github.com/AdRoll/goamz/s3"
	"github.com/stretchr/testify/assert"
)

func TestS3Mock(t *testing.T) {
	src := NewMock(aws.Auth{}, aws.Region{})
	bkt := src.Bucket("hello")

	// Write
	content := "Hello World!\n"
	assert.Nil(t,
		bkt.PutReader("hello.txt", strings.NewReader(content), int64(len(content)),
			"text/plain", s3.PublicRead, s3.Options{}))

	// Read
	rc, e := bkt.GetReader("hello.txt")
	assert.Nil(t, e)
	defer rc.Close()

	// Check equivalence
	b, e := ioutil.ReadAll(rc)
	assert.Nil(t, e)
	assert.Equal(t, content, string(b))

	// Overwrite
	assert.Nil(t,
		bkt.PutReader("hello.txt", strings.NewReader(content), int64(len(content)),
			"text/plain", s3.PublicRead, s3.Options{}))

	// Read and check equivalence again
	{
		rc, e := bkt.GetReader("hello.txt")
		assert.Nil(t, e)
		defer rc.Close()

		b, e := ioutil.ReadAll(rc)
		assert.Nil(t, e)
		assert.Equal(t, content, string(b))
	}
}
