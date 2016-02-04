package s3mock

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/AdRoll/goamz/aws"
	"github.com/AdRoll/goamz/s3"
	"github.com/stretchr/testify/assert"
)

func TestS3Mock(t *testing.T) {
	src := New(aws.Auth{}, aws.Region{})
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

func TestS3MockInUse(t *testing.T) {
	assert := assert.New(t)

	src := New(aws.Auth{}, aws.Region{})
	bkt := src.Bucket("hello")

	m := image.NewRGBA64(image.Rect(0, 0, 100, 100))
	writeImage(bkt, "a.jpg", m)

	assert.NotNil(bkt.Fs["a.jpg"])

	m1, _, e := readImage(bkt, "a.jpg")
	assert.Nil(e)
	assert.Equal(image.Rect(0, 0, 100, 100), m1.Bounds())

	assert.Nil(bkt.Fs["b.jpg"])
	_, _, e = readImage(bkt, "b.jpg")
	assert.NotNil(e)
}

type bucket interface {
	GetReader(p string) (io.ReadCloser, error)
	PutReader(p string, r io.Reader, length int64, contType string, perm s3.ACL, options s3.Options) error
}

func readImage(bkt bucket, filename string) (image.Image, string, error) {
	r, e := bkt.GetReader(filename)
	if e != nil {
		return nil, "", e
	}
	defer r.Close()
	return image.Decode(r)
}

// writeImage writes images in JPEG format.
func writeImage(bkt bucket, filename string, m image.Image) error {
	var buf bytes.Buffer
	if e := jpeg.Encode(&buf, m, &jpeg.Options{Quality: 7}); e != nil {
		return e
	}
	return bkt.PutReader(filename, &buf, int64(buf.Len()), "image/jpeg", s3.PublicRead, s3.Options{})
}
