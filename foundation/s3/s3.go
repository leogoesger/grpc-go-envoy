package s3

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Session struct ...
type Session struct {
	svc        *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

// New creates aws s3 session and creates the uploader and downloader
// along with session svc.
func New(region string) Session {
	// Initialize a session in region that the SDK will use to load
	// credentials from the .env file.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		exitErrorf("Unable to create session, %v", err)
	}

	// Create S3 service client, uploader and downloader.
	svc := s3.New(sess)
	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)

	return Session{svc, uploader, downloader}
}

// CreateBucket creates a bucket in the region configured in the shared config
// or AWS_REGION environment variable.
func (session *Session) CreateBucket(bucket string) {
	_, err := session.svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		exitErrorf("Unable to create bucket %q, %v", bucket, err)
	}

	err = session.svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		exitErrorf("Error occurred while waiting for bucket to be created, %v", bucket)
	}
	fmt.Printf("Bucket %q successfully created\n", bucket)
}

// ListBucketItems lists the items in the specified S3 Bucket
func (session *Session) ListBucketItems(bucket string) {
	resp, err := session.svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}

// UploadFile uploads the file's body to S3 bucket as an object with the key being the
// same as the filename.
func (session *Session) UploadFile(bucket string, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()

	_, err = session.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}

// DownloadFile an item from an S3 Bucket in the region configured in the shared config
// or AWS_REGION environment variable.
func (session *Session) DownloadFile(bucket string, item string) {
	file, err := os.Create(item)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()

	numBytes, err := session.downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", item, err)
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
