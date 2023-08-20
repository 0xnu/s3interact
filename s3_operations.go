package main

import (
	"fmt"
	"strings"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func createBucket(svc *s3.S3, bucket string) {
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		fmt.Println("Error creating bucket:", err)
		return
	}
	fmt.Println("Bucket created successfully.")
}

func createFolder(svc *s3.S3, bucket, folder string) {
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(folder + "/"),
	})
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}
	fmt.Println("Folder created successfully.")
}

func uploadSingleFile(svc *s3.S3, bucket, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
		Body:   file,
	})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}
	fmt.Println("File uploaded successfully.")
}

func uploadMultipleFiles(svc *s3.S3, bucket, filePaths string) {
	paths := strings.Split(filePaths, ",")
	for _, path := range paths {
		uploadSingleFile(svc, bucket, strings.TrimSpace(path))
	}
}

func deleteSingleFile(svc *s3.S3, bucket, fileKey string) {
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}
	fmt.Println("File deleted successfully.")
}

func deleteMultipleFiles(svc *s3.S3, bucket, fileKeys string) {
	keys := strings.Split(fileKeys, ",")
	objects := make([]*s3.ObjectIdentifier, len(keys))
	for i, key := range keys {
		objects[i] = &s3.ObjectIdentifier{Key: aws.String(strings.TrimSpace(key))}
	}

	_, err := svc.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{Objects: objects},
	})
	if err != nil {
		fmt.Println("Error deleting files:", err)
		return
	}
	fmt.Println("Files deleted successfully.")
}

func deleteFolder(svc *s3.S3, bucket, folder string) {
	// List objects in the folder
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(folder + "/"),
	})
	if err != nil {
		fmt.Println("Error listing objects:", err)
		return
	}

	// Delete objects and subfolders in the folder
	objects := make([]*s3.ObjectIdentifier, len(resp.Contents))
	for i, item := range resp.Contents {
		objects[i] = &s3.ObjectIdentifier{Key: item.Key}
	}

	if len(objects) > 0 {
		_, err = svc.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &s3.Delete{Objects: objects},
		})
		if err != nil {
			fmt.Println("Error deleting objects:", err)
			return
		}
	}

	// Delete the folder itself (represented as an object with a trailing slash)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(folder + "/"),
	})
	if err != nil {
		fmt.Println("Error deleting folder:", err)
		return
	}

	fmt.Println("Folder deleted successfully.")
}
