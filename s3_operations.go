package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

func createBucket(svc s3iface.S3API, bucket string) error {
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		fmt.Println("Error creating bucket:", err)
		return err
	}
	fmt.Println("Bucket created successfully.")
	return nil
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
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(folder + "/"),
	})
	if err != nil {
		fmt.Println("Error listing objects:", err)
		return
	}

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

func listBucketsAndObjects(svc *s3.S3) {
	result, err := svc.ListBuckets(nil)
	if err != nil {
		fmt.Println("Error listing buckets:", err)
		return
	}

	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("* %s\n", aws.StringValue(b.Name))

		resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: b.Name})
		if err != nil {
			fmt.Println("Error listing objects:", err)
			continue
		}

		fmt.Println("  Objects:")
		for _, item := range resp.Contents {
			fmt.Printf("    - %s\n", aws.StringValue(item.Key))
		}
	}
}

func downloadSingleFile(svc *s3.S3, bucket, fileKey, destinationPath string) {
	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer output.Body.Close()

	file, err := os.Create(destinationPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, output.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File downloaded successfully.")
}

func downloadMultipleFiles(svc *s3.S3, bucket string, fileKeysAndPaths map[string]string) {
	for fileKey, destinationPath := range fileKeysAndPaths {
		downloadSingleFile(svc, bucket, fileKey, destinationPath)
	}
}
