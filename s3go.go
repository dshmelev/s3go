package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	tags := flag.String("t", "", "tags for object")
	flag.Parse()
	filename := flag.Args()[0]
	bucket := flag.Args()[1]
	fmt.Printf("Upload %q to bucket %q with tag: %q\n", filename, bucket, *tags)

	// Initialize a session that the SDK will use to load configuration,
	// credentials, and region from the shared config file. (~/.aws/config).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := s3.New(sess)
	tag_set := parse_tags(tags)

	upload(&bucket, &filename, svc)
	tagging(&bucket, &filename, tag_set, svc)
}

func parse_tags(str *string) []*s3.Tag {
	// Parse tags string to []*s3.Tag
	var tags []*s3.Tag
	for _, element := range strings.Split(*str, ",") {
		pair := strings.Split(element, "=")
		tags = append(tags, &s3.Tag{
			Key:   &pair[0],
			Value: &pair[1],
		})
	}
	return tags
}

func upload(bucket *string, filename *string, svc *s3.S3) {
	// Try to open the file
	file, err := os.Open(*filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()
	// Upload S3 object
	params := &s3.PutObjectInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*filename),
		Body:   file,
	}
	_, err = svc.PutObject(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func tagging(bucket *string, filename *string, tag_set []*s3.Tag, svc *s3.S3) {
	// Tagging S3 object
	params := &s3.PutObjectTaggingInput{
		Bucket: aws.String(*bucket),   // Required
		Key:    aws.String(*filename), // Required
		Tagging: &s3.Tagging{ // Required
			TagSet: tag_set,
		},
	}
	_, err := svc.PutObjectTagging(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
