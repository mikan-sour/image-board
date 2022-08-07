package main

import (
	"fmt"
	"image-board/src/batch/s3Client"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func main() {

	myBucket := "saint-bernards"
	myPups := []string{"./images/sb1.jpg", "./images/sb2.jpg"}

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

	// check if bucket "saint-bernards" exists
	var exists = false
	buckets, err := client.ListBuckets()
	if err != nil {
		panic(err)
	}
	for _, b := range buckets {
		if *b.Name == myBucket {
			exists = true
			break
		}
	}

	// if does not exist, create bucket
	if !exists {
		err = client.CreateBucket(myBucket)
		if err != nil {
			panic(err)
		}
	}

	// add a file to the bucket
	for _, file := range myPups {
		err := client.PutObject(myBucket, file)
		if err != nil {
			panic(err)
		}
	}

	// check that files exist
	filesInBucket, err := client.ListItems(myBucket)
	if err != nil {
		panic(err)
	}
	for i, file := range filesInBucket {
		fmt.Printf("%d:%s\n", i, *file.Key)
	}

	// download one
	// file, err := client.GetObject(myBucket, myPups[0])
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("got file %s\n", file.Name())

	// // delete
	// for _, file := range myPups {
	// 	err := client.DeleteObject(myBucket, file)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	fmt.Println("end of the road")
	os.Exit(0)
}
