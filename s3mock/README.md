# s3mock

s3mock provides an interface type `s3mock.MockableBucket` which
abstracts both `github.com/AdRoll/goamz/s3.Bucket` and
`s3mock.bucket`.

(Currently `s3mock.MockableBucket` doesn't yet implement all methods
of `github.com/AdRoll/goamz/s3.Bucket`.  But that is our future goal.)

Like we create `s3.Bucket` by calling `s3.New().Bucket(name)`, we can
create `s3mock.bucket` by calling `s3mock.New().Bucket(name)`.

To write testable Go code that accesses AWS S3 service, please make it
uses `s3mock.MockableBucket`, instead of `s3.Bucket`, as the interface
to S3.  Then we can test the code by giving it a mock bucket from
`s3mock.New().Bucket(name)`.  And we can use it with real AWS S3 by
giving it a real bucket from `s3.New().Bucket(name)`.

For more details, please refer to https://github.com/topicai/goamzmock/blob/master/s3mock/example/images.go
