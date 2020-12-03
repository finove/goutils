package s3file

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Store access aws s3
type S3Store struct {
	endpoint string
	region   string
	sess     *session.Session
	svc      *s3.S3
}

// KeyInfo file status
type KeyInfo struct {
	Key   string    `json:"key,omitempty"`
	Size  int64     `json:"size,omitempty"`
	Stamp time.Time `json:"stamp,omitempty"`
}

// Connect connect to s3 server
func (s *S3Store) Connect(accessID, secretKey string, region string, endpoints ...string) (err error) {
	if len(endpoints) > 0 {
		s.endpoint = endpoints[0]
	}
	s.region = region
	s.sess, err = session.NewSession(&aws.Config{
		Endpoint:         aws.String(s.endpoint),
		Region:           aws.String(s.region),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessID, secretKey, ""),
	})
	if err != nil {
		return
	}
	s.svc = s3.New(s.sess)
	return
}

// ListBuckets Returns a list of all buckets owned
func (s *S3Store) ListBuckets() (bkts []string, err error) {
	var req s3.ListBucketsInput
	var resp *s3.ListBucketsOutput
	if resp, err = s.svc.ListBuckets(&req); err != nil {
		return
	}
	for _, bkt := range resp.Buckets {
		bkts = append(bkts, aws.StringValue(bkt.Name))
	}
	return
}

// ListBucketFiles Returns some or all (up to 1,000) of the objects in a bucket
func (s *S3Store) ListBucketFiles(bucket, prefix string) (prefixs []string, keys []KeyInfo, err error) {
	var req s3.ListObjectsInput
	var resp *s3.ListObjectsOutput
	req.SetBucket(bucket)
	req.SetPrefix(prefix)
	req.SetDelimiter("/")
	if resp, err = s.svc.ListObjects(&req); err != nil {
		return
	}
	for _, c := range resp.Contents {
		var info KeyInfo
		info.Key = aws.StringValue(c.Key)
		info.Size = aws.Int64Value(c.Size)
		info.Stamp = aws.TimeValue(c.LastModified)
		keys = append(keys, info)
	}
	for _, c := range resp.CommonPrefixes {
		prefixs = append(prefixs, aws.StringValue(c.Prefix))
	}
	return
}

// GetPresignedURL creates a presigned URL for a bucket object
func (s *S3Store) GetPresignedURL(bucket, key string, durations ...time.Duration) (string, error) {
	var expire time.Duration
	if len(durations) > 0 {
		expire = durations[0]
	} else {
		expire = 15 * time.Minute
	}
	req, _ := s.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(expire)
	if err != nil {
		return "", err
	}
	return urlStr, nil
}
