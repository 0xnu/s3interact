package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
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

func deleteBucket(svc s3iface.S3API, region string, bucket string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return err
	}

	newS3Client := s3.New(sess)

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}

	_, err = newS3Client.DeleteBucket(input)
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

func setRegion(svc *s3.S3, region string) {
	svc.Config.Region = aws.String(region)
	fmt.Println("Region set successfully to:", region)
}

func moveFiles(svc *s3.S3, bucket, sourceFolder, destinationFolder string, fileKeys []string) {
	for _, fileKey := range fileKeys {
		sourceKey := sourceFolder + "/" + fileKey
		destinationKey := destinationFolder + "/" + fileKey

		_, err := svc.CopyObject(&s3.CopyObjectInput{
			Bucket:     aws.String(bucket),
			CopySource: aws.String(bucket + "/" + sourceKey),
			Key:        aws.String(destinationKey),
		})
		if err != nil {
			fmt.Println("Error copying file:", err)
			continue
		}

		_, err = svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(sourceKey),
		})
		if err != nil {
			fmt.Println("Error deleting original file:", err)
			continue
		}

		fmt.Printf("File %s moved successfully from %s to %s.\n", fileKey, sourceFolder, destinationFolder)
	}
}

func renameFile(svc *s3.S3, bucket, originalKey, newKey string) {
	_, err := svc.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(bucket + "/" + originalKey),
		Key:        aws.String(newKey),
	})
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(originalKey),
	})
	if err != nil {
		fmt.Println("Error deleting original file:", err)
		return
	}

	fmt.Printf("File %s renamed successfully to %s.\n", originalKey, newKey)
}

func moveFolders(svc *s3.S3, bucket string, sourceFolders, destinationFolders []string) {
	for i, sourceFolder := range sourceFolders {
		destinationFolder := destinationFolders[i]

		resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(sourceFolder + "/"),
		})
		if err != nil {
			fmt.Println("Error listing objects:", err)
			continue
		}

		for _, item := range resp.Contents {
			sourceKey := aws.StringValue(item.Key)
			destinationKey := strings.Replace(sourceKey, sourceFolder, destinationFolder, 1)

			_, err := svc.CopyObject(&s3.CopyObjectInput{
				Bucket:     aws.String(bucket),
				CopySource: aws.String(bucket + "/" + sourceKey),
				Key:        aws.String(destinationKey),
			})
			if err != nil {
				fmt.Println("Error copying object:", err)
				continue
			}

			_, err = svc.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(sourceKey),
			})
			if err != nil {
				fmt.Println("Error deleting original object:", err)
				continue
			}
		}

		fmt.Printf("Folder %s moved successfully to %s.\n", sourceFolder, destinationFolder)
	}
}

func renameFolders(svc *s3.S3, bucket string, originalFolders, newFolders []string) {
	for i, originalFolder := range originalFolders {
		newFolder := newFolders[i]

		resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(originalFolder + "/"),
		})
		if err != nil {
			fmt.Println("Error listing objects:", err)
			continue
		}

		for _, item := range resp.Contents {
			originalKey := aws.StringValue(item.Key)
			newKey := strings.Replace(originalKey, originalFolder, newFolder, 1)

			_, err := svc.CopyObject(&s3.CopyObjectInput{
				Bucket:     aws.String(bucket),
				CopySource: aws.String(bucket + "/" + originalKey),
				Key:        aws.String(newKey),
			})
			if err != nil {
				fmt.Println("Error copying object:", err)
				continue
			}

			_, err = svc.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(originalKey),
			})
			if err != nil {
				fmt.Println("Error deleting original object:", err)
				continue
			}
		}

		fmt.Printf("Folder %s renamed successfully to %s.\n", originalFolder, newFolder)
	}
}
