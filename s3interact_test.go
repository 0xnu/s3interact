package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockS3Client struct {
	s3iface.S3API
	buckets map[string]bool
}

func (m *mockS3Client) CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	if m.buckets[*input.Bucket] {
		return nil, awserr.New(s3.ErrCodeBucketAlreadyExists, "Bucket already exists", nil)
	}
	m.buckets[*input.Bucket] = true
	return &s3.CreateBucketOutput{}, nil
}

func (m *mockS3Client) DeleteBucket(input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	if !m.buckets[*input.Bucket] {
		return nil, awserr.New(s3.ErrCodeNoSuchBucket, "Bucket does not exist", nil)
	}
	delete(m.buckets, *input.Bucket)
	return &s3.DeleteBucketOutput{}, nil
}

func TestCreateAndDeleteBucket(t *testing.T) {
	svc := &mockS3Client{buckets: make(map[string]bool)}
	bucket := "s3interact-demo"

	err := createBucket(svc, bucket)
	if err != nil {
		t.Errorf("Expected no error in creating bucket, got %v", err)
	}

	err = deleteBucket(svc, bucket)
	if err != nil {
		t.Errorf("Expected no error in deleting bucket, got %v", err)
	}

	err = deleteBucket(svc, "non-existing-bucket")
	if err == nil {
		t.Errorf("Expected error for deleting non-existing bucket, got nil")
	}
}
