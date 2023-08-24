## S3Interact

[![Release](https://img.shields.io/github/release/0xnu/s3interact.svg)](https://github.com/0xnu/s3interact/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xnu/s3interact)](https://goreportcard.com/report/github.com/0xnu/s3interact)
[![Go Reference](https://pkg.go.dev/badge/github.com/0xnu/s3interact.svg)](https://pkg.go.dev/github.com/0xnu/s3interact)
[![License](https://img.shields.io/github/license/0xnu/s3interact)](/LICENSE)

S3interact provides a command-line interface for interacting with Amazon S3, enabling users to manage buckets, folders, and files easily. Users can create and delete buckets, folders, and files and upload single or multiple files through simple prompts and inputs, making it a resourceful tool for anyone working with Amazon S3.

Now users can set and delete bucket policies and set ACLs. You can set predefined ACL to values like `private`, `public-read`, `public-read-write`, `authenticated-read`, `aws-exec-read`, `bucket-owner-read`, `bucket-owner-full-control` and `log-delivery-write`, as per the [AWS S3 documentation](https://docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html).

Please carefully handle the policy JSON string when setting a bucket policy, as incorrect JSON or permissions can lead to unexpected behaviour or access issues.

### Execute Locally

Run the command in your terminal to execute the code.

```sh
make execute
```

### Build

Build single binary for local os.

```sh
make build
```

Build for multi os (linux 386, amd64).

```sh
make package
```

### To Do

- [x] ~~Recursive File/Folder Deletion~~
- [x] ~~List Buckets and Objects~~
- [x] ~~Downloading Files~~
- [x] ~~Bucket and Object Information~~
- [x] ~~Bucket Policies and Permissions~~
- [x] ~~Bucket Deletion~~
- [x] ~~Set a new Region~~
- [x] ~~Moving and Renaming Files~~
- [x] ~~Moving and Renaming Folders~~
- [x] ~~Generate Pre-Signed URL for an Object~~

### Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

### Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/Cloudeya/coronavirusapi-wrapper/tags).

### License

This project is licensed under the [BSD 3-Clause License](./LICENSE).

### Copyright

(c) 2023 [Finbarrs Oketunji](https://finbarrs.eu).