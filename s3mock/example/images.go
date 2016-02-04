package example

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/AdRoll/goamz/s3"
	"github.com/topicai/goamzmock/s3mock"
)

func readImage(bkt s3mock.MockableBucket, filename string) (image.Image, string, error) {
	r, e := bkt.GetReader(filename)
	if e != nil {
		return nil, "", e
	}
	defer r.Close()
	return image.Decode(r)
}

func writeImage(bkt s3mock.MockableBucket, filename string, m image.Image) error {
	var buf bytes.Buffer
	if e := jpeg.Encode(&buf, m, &jpeg.Options{Quality: 7}); e != nil {
		return e
	}
	return bkt.PutReader(filename, &buf, int64(buf.Len()), "image/jpeg", s3.PublicRead, s3.Options{})
}
