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
		fmt.Println("7. Download a single file")
		fmt.Println("8. Download multiple files")
		fmt.Println("9. List Buckets and Objects")
		fmt.Println("10. Exit")
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
			fmt.Print("Enter file key: ")
			fileKey, _ := reader.ReadString('\n')
			fmt.Print("Enter destination path: ")
			destinationPath, _ := reader.ReadString('\n')
			downloadSingleFile(svc, bucket, strings.TrimSpace(fileKey), strings.TrimSpace(destinationPath))
		case "8":
			fmt.Print("Enter file keys and destination paths (comma-separated, key:path): ")
			fileKeysAndPathsInput, _ := reader.ReadString('\n')
			fileKeysAndPaths := make(map[string]string)
			pairs := strings.Split(fileKeysAndPathsInput, ",")
			for _, pair := range pairs {
				keyAndPath := strings.Split(strings.TrimSpace(pair), ":")
				if len(keyAndPath) == 2 {
					fileKeysAndPaths[keyAndPath[0]] = keyAndPath[1]
				}
			}
			downloadMultipleFiles(svc, bucket, fileKeysAndPaths)
		case "9":
			listBucketsAndObjects(svc)
		case "10":
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
