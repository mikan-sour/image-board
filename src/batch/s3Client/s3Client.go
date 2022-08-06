package s3Client

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Client interface {
	Create(bucketName string) error
	List() ([]*s3.Bucket, error)
	Put(bucketName, fileName string) error
	Get(name string)
	ListItems(bucketName string)
	Delete(name string)
}

type S3ClientImpl struct {
	svc     *s3.S3
	session *session.Session
}

func NewClient(config *aws.Config) (*S3ClientImpl, error) {
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &S3ClientImpl{
		svc:     s3.New(sess),
		session: sess,
	}, nil
}

func (s *S3ClientImpl) Create(bucketName string) error {
	_, err := s.svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}

	fmt.Printf("Waiting for bucket %q to be created...\n", bucketName)

	err = s.svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}

	fmt.Printf("Bucket %q successfully created\n", bucketName)
	return nil
}

func (s *S3ClientImpl) List() ([]*s3.Bucket, error) {
	var buckets = []*s3.Bucket{}
	result, err := s.svc.ListBuckets(nil)
	if err != nil {
		return buckets, err
	}

	buckets = append(buckets, result.Buckets...)

	return buckets, nil
}

func (s *S3ClientImpl) Put(bucketName, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	uploader := s3manager.NewUploader(s.session)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucketName)
	return nil
}
