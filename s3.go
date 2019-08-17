package hemmingway


import (
	"encoding/csv"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"strings"
)

// Downloads a File from S3 and returns a Reader pointer.
// The Reader pointer can then be used by .Read() or .ReadAll()
// Credentials are expected to be stored in a AWS_ACCESS_KEY and AWS_SECRET_KEY env variable
func FetchS3(bucket string, filename string, region string) *csv.Reader {
	sess, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewEnvCredentials(),
		Region: aws.String(region),
	},
	)

	downloader := s3manager.NewDownloader(sess)

	file := &aws.WriteAtBuffer{}
	_, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		})

	if err != nil {
		log.Fatalf("Unable to download item %q, %v", filename, err)
	}

	csvReader := csv.NewReader(strings.NewReader(string(file.Bytes())))

	return csvReader

}

