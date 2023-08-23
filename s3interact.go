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

	actions := map[string]func(*s3.S3, string, *bufio.Reader){
		"1":  createFolderAction,
		"2":  uploadSingleFileAction,
		"3":  uploadMultipleFilesAction,
		"4":  deleteSingleFileAction,
		"5":  deleteMultipleFilesAction,
		"6":  deleteFolderAction,
		"7":  downloadSingleFileAction,
		"8":  downloadMultipleFilesAction,
		"9":  listBucketsAndObjectsAction,
		"10": getBucketInfoAction,
		"11": getObjectInfoAction,
		"12": setBucketPolicyAction,
		"13": deleteBucketPolicyAction,
		"14": setBucketACLAction,
		"15": deleteBucketAction,
		"16": setRegionAction,
		"17": moveFilesAction,
		"18": renameFileAction,
		"19": moveFoldersAction,
		"20": renameFoldersAction,
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
		fmt.Println("10. Get Bucket Information")
		fmt.Println("11. Get Object Information")
		fmt.Println("12. Set Bucket Policy")
		fmt.Println("13. Delete Bucket Policy")
		fmt.Println("14. Set Bucket ACL")
		fmt.Println("15. Delete Bucket")
		fmt.Println("16. Set a Region")
		fmt.Println("17. Move a File")
		fmt.Println("18. Rename a File")
		fmt.Println("19. Move a Folder")
		fmt.Println("20. Rename a Folder")
		fmt.Println("21. Exit")
		fmt.Print("Enter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		action, exists := actions[choice]
		if exists {
			action(svc, bucket, reader)
		} else if choice == "20" {
			return
		} else {
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func createFolderAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter folder name: ")
	folder, _ := reader.ReadString('\n')
	createFolder(svc, bucket, strings.TrimSpace(folder))
}

func uploadSingleFileAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter file path: ")
	filePath, _ := reader.ReadString('\n')
	uploadSingleFile(svc, bucket, strings.TrimSpace(filePath))
}

func uploadMultipleFilesAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter file paths (comma-separated): ")
	filePaths, _ := reader.ReadString('\n')
	uploadMultipleFiles(svc, bucket, strings.TrimSpace(filePaths))
}

func deleteSingleFileAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter file key: ")
	fileKey, _ := reader.ReadString('\n')
	deleteSingleFile(svc, bucket, strings.TrimSpace(fileKey))
}

func deleteMultipleFilesAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter file keys (comma-separated): ")
	fileKeys, _ := reader.ReadString('\n')
	deleteMultipleFiles(svc, bucket, strings.TrimSpace(fileKeys))
}

func deleteFolderAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter folder name: ")
	folder, _ := reader.ReadString('\n')
	deleteFolder(svc, bucket, strings.TrimSpace(folder))
}

func downloadSingleFileAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter file key: ")
	fileKey, _ := reader.ReadString('\n')
	fmt.Print("Enter destination path: ")
	destinationPath, _ := reader.ReadString('\n')
	downloadSingleFile(svc, bucket, strings.TrimSpace(fileKey), strings.TrimSpace(destinationPath))
}

func downloadMultipleFilesAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
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
}

func listBucketsAndObjectsAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	listBucketsAndObjects(svc)
}

func getBucketInfoAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter bucket name: ")
	bucketName, _ := reader.ReadString('\n')
	getBucketInfo(svc, strings.TrimSpace(bucketName))
}

func getObjectInfoAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter bucket name: ")
	bucketName, _ := reader.ReadString('\n')
	fmt.Print("Enter object key: ")
	objectKey, _ := reader.ReadString('\n')
	getObjectInfo(svc, strings.TrimSpace(bucketName), strings.TrimSpace(objectKey))
}

func setBucketPolicyAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter bucket name: ")
	bucketName, _ := reader.ReadString('\n')
	fmt.Print("Enter policy JSON: ")
	policy, _ := reader.ReadString('\n')
	setBucketPolicy(svc, strings.TrimSpace(bucketName), strings.TrimSpace(policy))
}

func deleteBucketPolicyAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter bucket name: ")
	bucketName, _ := reader.ReadString('\n')
	deleteBucketPolicy(svc, strings.TrimSpace(bucketName))
}

func setBucketACLAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter bucket name: ")
	bucketName, _ := reader.ReadString('\n')
	fmt.Print("Enter ACL (e.g., private, public-read): ")
	acl, _ := reader.ReadString('\n')
	setBucketACL(svc, strings.TrimSpace(bucketName), strings.TrimSpace(acl))
}

func deleteBucketAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter bucket name: ")
	bucketName, _ := reader.ReadString('\n')
	bucketName = strings.TrimSpace(bucketName)
	region := *svc.Config.Region
	deleteBucket(svc, region, bucketName)
}

func setRegionAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter new AWS Region (e.g., eu-west-2): ")
	newRegion, _ := reader.ReadString('\n')
	newRegion = strings.TrimSpace(newRegion)
	setRegion(svc, newRegion)
}

func moveFilesAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter source folder: ")
	sourceFolder, _ := reader.ReadString('\n')
	sourceFolder = strings.TrimSpace(sourceFolder)

	fmt.Print("Enter destination folder: ")
	destinationFolder, _ := reader.ReadString('\n')
	destinationFolder = strings.TrimSpace(destinationFolder)

	fmt.Print("Enter file keys to move (comma-separated): ")
	fileKeysInput, _ := reader.ReadString('\n')
	fileKeys := strings.Split(strings.TrimSpace(fileKeysInput), ",")

	for i, key := range fileKeys {
		fileKeys[i] = strings.TrimSpace(key)
	}

	moveFiles(svc, bucket, sourceFolder, destinationFolder, fileKeys)
}

func renameFileAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter original file key: ")
	originalKey, _ := reader.ReadString('\n')
	originalKey = strings.TrimSpace(originalKey)

	fmt.Print("Enter new file key: ")
	newKey, _ := reader.ReadString('\n')
	newKey = strings.TrimSpace(newKey)

	renameFile(svc, bucket, originalKey, newKey)
}

func moveFoldersAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter source folders (comma-separated): ")
	sourceFoldersInput, _ := reader.ReadString('\n')
	sourceFolders := strings.Split(strings.TrimSpace(sourceFoldersInput), ",")

	fmt.Print("Enter destination folders (comma-separated): ")
	destinationFoldersInput, _ := reader.ReadString('\n')
	destinationFolders := strings.Split(strings.TrimSpace(destinationFoldersInput), ",")

	if len(sourceFolders) != len(destinationFolders) {
		fmt.Println("Error: The number of source folders must match the number of destination folders.")
		return
	}

	for i, folder := range sourceFolders {
		sourceFolders[i] = strings.TrimSpace(folder)
		destinationFolders[i] = strings.TrimSpace(destinationFolders[i])
	}

	moveFolders(svc, bucket, sourceFolders, destinationFolders)
}

func renameFoldersAction(svc *s3.S3, bucket string, reader *bufio.Reader) {
	fmt.Print("Enter original folder names (comma-separated): ")
	originalFoldersInput, _ := reader.ReadString('\n')
	originalFolders := strings.Split(strings.TrimSpace(originalFoldersInput), ",")

	fmt.Print("Enter new folder names (comma-separated): ")
	newFoldersInput, _ := reader.ReadString('\n')
	newFolders := strings.Split(strings.TrimSpace(newFoldersInput), ",")

	if len(originalFolders) != len(newFolders) {
		fmt.Println("Error: The number of original folders must match the number of new folder names.")
		return
	}

	for i, folder := range originalFolders {
		originalFolders[i] = strings.TrimSpace(folder)
		newFolders[i] = strings.TrimSpace(newFolders[i])
	}

	renameFolders(svc, bucket, originalFolders, newFolders)
}
