package example

import (
	"image"
	"testing"

	"github.com/AdRoll/goamz/aws"
	"github.com/stretchr/testify/assert"
	"github.com/topicai/goamzmock/s3mock"
)

func TestS3MockInUse(t *testing.T) {
	assert := assert.New(t)

	src := s3mock.NewMock(aws.Auth{}, aws.Region{})
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
