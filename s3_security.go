package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getBucketInfo(svc *s3.S3, bucket string) {
	input := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.GetBucketLocation(input)
	if err != nil {
		fmt.Println("Error getting bucket information:", err)
		return
	}

	fmt.Printf("Bucket: %s\n", bucket)
	fmt.Printf("Location: %s\n", aws.StringValue(result.LocationConstraint))
}

func getObjectInfo(svc *s3.S3, bucket, objectKey string) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		fmt.Println("Error getting object information:", err)
		return
	}

	fmt.Printf("Object Key: %s\n", objectKey)
	fmt.Printf("Size: %d bytes\n", aws.Int64Value(result.ContentLength))
	fmt.Printf("Last Modified: %s\n", aws.TimeValue(result.LastModified))
	fmt.Printf("Content Type: %s\n", aws.StringValue(result.ContentType))
}
