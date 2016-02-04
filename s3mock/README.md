# s3mock

s3mock provides a mock implementation of Go package
github.com/AdRoll/goamz/s3.  In particular, it mocks two types:

1. S3, and
1. Bucket

It is notable that currently s3mock.Bucket doesn't yet implement all
methods of s3.Bucket.  But that is our future goal.
