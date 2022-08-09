package batch

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jedzeins/image-board/src/s3Client"

	"github.com/aws/aws-sdk-go/aws"
)

type Batch interface {
	CheckLocalStack() error
	GetFilePaths(dir string) ([]string, error)
	PreloadImages(dirs []string) error
}

type BatchImpl struct {
	bucketName  string
	svc         s3Client.S3Client
	healthcheck Healthcheck
}

func NewBatch(bucketName string, bucketURL string, s3Config *aws.Config) (*BatchImpl, error) {
	svc, err := s3Client.NewClient(s3Config)
	if err != nil {
		return nil, err
	}

	return &BatchImpl{
		bucketName:  bucketName,
		svc:         svc,
		healthcheck: NewHealthcheck(bucketURL),
	}, nil
}

func (b *BatchImpl) CheckLocalStack() error {
	err := b.healthcheck.DoWhile()
	if err != nil {
		log.Fatalf("BatchImpl.checkLocalStack - CheckLocalStack: %v\n", err)
	}
	return nil
}

func (b *BatchImpl) GetFilePaths(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("BatchImpl.getFilePaths: %v\n", err)
		return nil, err
	}

	var namesOfFiles []string
	for _, f := range files {
		namesOfFiles = append(namesOfFiles, fmt.Sprintf("%s/%s", dir, f.Name()))
	}

	return namesOfFiles, nil
}

func (b *BatchImpl) PreloadImages(dirs []string) error {
	var exists = false
	buckets, err := b.svc.ListBuckets()
	if err != nil {
		log.Fatalf("BatchImpl.preloadImages - ListBuckets: %v\n", err)
	}
	for _, bucket := range buckets {
		if *bucket.Name == b.bucketName {
			exists = true
			break
		}
	}

	// if does not exist, create bucket
	if !exists {
		err = b.svc.CreateBucket(b.bucketName)
		if err != nil {
			log.Fatalf("BatchImpl.preloadImages - CreateBucket: %v\n", err)
		}
	}

	// add a file to the bucket
	for _, file := range dirs {
		err := b.svc.PutObject(b.bucketName, file)
		if err != nil {
			log.Fatalf("BatchImpl.preloadImages - PutObject: %v\n", err)
		}
	}

	// check that files exist
	_, err = b.svc.ListItems(b.bucketName)
	if err != nil {
		log.Fatalf("BatchImpl.preloadImages - ListItems: %v\n", err)
	}

	fmt.Println("load images to s3 successfully executed")
	return nil
}
