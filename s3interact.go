package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter AWS Key ID: ")
	awsKeyID, _ := reader.ReadString('\n')

	fmt.Print("Enter AWS Secret Key: ")
	awsSecretKey, _ := reader.ReadString('\n')

	fmt.Print("Enter AWS Region (e.g., eu-west-2): ")
	awsRegion, _ := reader.ReadString('\n')

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(strings.TrimSpace(awsRegion)),
		Credentials: credentials.NewStaticCredentials(strings.TrimSpace(awsKeyID), strings.TrimSpace(awsSecretKey), ""),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	svc := s3.New(sess)

	fmt.Print("Do you want to create a new bucket? (yes/no): ")
	createBucketChoice, _ := reader.ReadString('\n')
	createBucketChoice = strings.TrimSpace(createBucketChoice)

	var bucket string
	if createBucketChoice == "yes" {
		fmt.Print("Enter new bucket name: ")
		bucket, _ = reader.ReadString('\n')
		bucket = strings.TrimSpace(bucket)
		createBucket(svc, bucket)
	} else {
		fmt.Print("Enter existing bucket name: ")
		bucket, _ = reader.ReadString('\n')
		bucket = strings.TrimSpace(bucket)
	}

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Create a folder")
		fmt.Println("2. Upload a single file")
		fmt.Println("3. Upload multiple files")
		fmt.Println("4. Delete a single file")
		fmt.Println("5. Delete multiple files")
		fmt.Println("6. Delete a folder")
		fmt.Println("7. Exit")
		fmt.Print("Enter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter folder name: ")
			folder, _ := reader.ReadString('\n')
			createFolder(svc, bucket, strings.TrimSpace(folder))
		case "2":
			fmt.Print("Enter file path: ")
			filePath, _ := reader.ReadString('\n')
			uploadSingleFile(svc, bucket, strings.TrimSpace(filePath))
		case "3":
			fmt.Print("Enter file paths (comma-separated): ")
			filePaths, _ := reader.ReadString('\n')
			uploadMultipleFiles(svc, bucket, strings.TrimSpace(filePaths))
		case "4":
			fmt.Print("Enter file key: ")
			fileKey, _ := reader.ReadString('\n')
			deleteSingleFile(svc, bucket, strings.TrimSpace(fileKey))
		case "5":
			fmt.Print("Enter file keys (comma-separated): ")
			fileKeys, _ := reader.ReadString('\n')
			deleteMultipleFiles(svc, bucket, strings.TrimSpace(fileKeys))
		case "6":
			fmt.Print("Enter folder name: ")
			folder, _ := reader.ReadString('\n')
			deleteFolder(svc, bucket, strings.TrimSpace(folder))
		case "7":
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

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

	// Delete objects in the folder
	objects := make([]*s3.ObjectIdentifier, len(resp.Contents))
	for i, item := range resp.Contents {
		objects[i] = &s3.ObjectIdentifier{Key: item.Key}
	}

	_, err = svc.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{Objects: objects},
	})
	if err != nil {
		fmt.Println("Error deleting folder:", err)
		return
	}
	fmt.Println("Folder deleted successfully.")
}
