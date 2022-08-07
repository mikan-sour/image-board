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
	CreateBucket(bucketName string) error
	ListBuckets() ([]*s3.Bucket, error)
	PutObject(bucketName, fileName string) error
	GetObject(bucketName, fileName string) (*os.File, error)
	ListItems(bucketName string) (*[]*s3.Object, error)
	DeleteObject(name string)
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

func (s *S3ClientImpl) CreateBucket(bucketName string) error {
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

func (s *S3ClientImpl) ListBuckets() ([]*s3.Bucket, error) {
	result, err := s.svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	return result.Buckets, nil
}

func (s *S3ClientImpl) PutObject(bucketName, fileName string) error {
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

func (s *S3ClientImpl) GetObject(bucketName, fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(s.session)

	_, err = downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileName),
		},
	)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *S3ClientImpl) ListItems(bucketName string) ([]*s3.Object, error) {
	resp, err := s.svc.ListObjectsV2(
		&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)},
	)
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3ClientImpl) DeleteObject(bucketName, fileName string) error {
	_, err := s.svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(fileName)})
	if err != nil {
		return err
	}

	err = s.svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return err
	}
	fmt.Printf("Object %s successfully deleted from bucket %s", fileName, bucketName)
	return nil
}
