package main

import (
	"fmt"
	"image-board/src/batch/s3Client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func main() {

	config := &aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4566"),
	}

	client, err := s3Client.NewClient(config)
	if err != nil {
		panic(err)
	}

	err = client.Put("dogs", "./src/batch/images/sb2.jpg")
	if err != nil {
		panic(err)
	}

	fmt.Println("put something?")
}
