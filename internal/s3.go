package internal

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
)

type S3Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

type NewRecordRequest struct {
	Key   string
	Value string
}

type GetRecordRequest struct {
	Key string
}

func deleteAllObjects(objList *s3.ListObjectsOutput, svc *s3.S3, config S3Config) error {
	for _, obj := range objList.Contents {
		_, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(config.BucketName),
			Key:    obj.Key,
		})
		if err != nil {
			return fmt.Errorf("failed to delete object: %s %v", *obj.Key, err)
		}
	}
	fmt.Println("All objects deleted successfully!")
	return nil
}

func listAllObjects(svc *s3.S3, config S3Config) (*s3.ListObjectsOutput, error) {
	return svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(config.BucketName),
	})
}

func newSession(config S3Config) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	})
	return sess, err
}

func newS3Service(config S3Config) (*s3.S3, error) {
	sess, err := newSession(config)
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}

func DropDB(config S3Config) error {
	svc, err := newS3Service(config)
	if err != nil {
		return err
	}

	allObjectsList, err := listAllObjects(svc, config)
	if err != nil {
		return err
	}

	if err := deleteAllObjects(allObjectsList, svc, config); err != nil {
		return err
	}

	return nil
}

func NewRecord(config S3Config, request NewRecordRequest) error {
	svc, err := newS3Service(config)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(request.Key),
		Body:   newBody(request.Value),
	}

	_, err = svc.PutObject(input)
	if err != nil {
		return err
	}

	fmt.Println("Object added successfully!")
	return nil
}

func GetRecord(config S3Config, request GetRecordRequest) (string, error) {
	svc, err := newS3Service(config)
	if err != nil {
		return "", err
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(request.Key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return "object not found", nil
		}
		return "", err
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(result.Body); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func newBody(value string) io.ReadSeeker {
	return bytes.NewReader([]byte(value))
}
