package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
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

func setBucketPolicy(svc *s3.S3, bucket, policy string) {
	input := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(policy),
	}

	_, err := svc.PutBucketPolicy(input)
	if err != nil {
		fmt.Println("Error setting bucket policy:", err)
		return
	}

	fmt.Println("Bucket policy set successfully.")
}

func deleteBucketPolicy(svc *s3.S3, bucket string) {
	input := &s3.DeleteBucketPolicyInput{
		Bucket: aws.String(bucket),
	}

	_, err := svc.DeleteBucketPolicy(input)
	if err != nil {
		fmt.Println("Error deleting bucket policy:", err)
		return
	}

	fmt.Println("Bucket policy deleted successfully.")
}

func setBucketACL(svc *s3.S3, bucket, acl string) {
	validACLs := []string{"private", "public-read", "public-read-write", "authenticated-read", "aws-exec-read", "bucket-owner-read", "bucket-owner-full-control", "log-delivery-write"}
	isValidACL := false
	for _, validACL := range validACLs {
		if acl == validACL {
			isValidACL = true
			break
		}
	}

	if !isValidACL {
		fmt.Println("Invalid ACL value. Please use one of the following:", strings.Join(validACLs, ", "))
		return
	}

	input := &s3.PutBucketAclInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}

	_, err := svc.PutBucketAcl(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println("Error Code:", aerr.Code())
			fmt.Println("Message:", aerr.Message())
			fmt.Println("Orig Err:", aerr.OrigErr())
		} else {
			fmt.Println(err.Error())
		}
		fmt.Println("Error setting bucket ACL:", err)
		return
	}

	fmt.Println("Bucket ACL set successfully.")
}

func deleteBucket(svc s3iface.S3API, bucket string) error {
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}

	_, err := svc.DeleteBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println("Error Code:", aerr.Code())
			fmt.Println("Message:", aerr.Message())
			fmt.Println("Orig Err:", aerr.OrigErr())
		} else {
			fmt.Println(err.Error())
		}
		fmt.Println("Error deleting bucket:", err)
		return err
	}

	fmt.Println("Bucket deleted successfully.")
	return nil
}
