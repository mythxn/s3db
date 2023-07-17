package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
)

type S3Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
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

func NewRecord(config S3Config, key string, value string) error {
	svc, err := newS3Service(config)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(key),
		Body:   strings.NewReader(value),
	}

	_, err = svc.PutObject(input)
	if err != nil {
		return err
	}

	fmt.Println("Object added successfully!")
	return nil
}

func GetRecord(config S3Config, key string) (string, error) {
	svc, err := newS3Service(config)
	if err != nil {
		return "", err
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		if IsObjectNotFoundErr(err) {
			return "", errors.New("object not found")
		}

	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(result.Body); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func IsObjectNotFoundErr(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == s3.ErrCodeNoSuchKey {
			return true
		}
	}
	return false
}

func ListAllObjects(config S3Config) ([]string, error) {
	svc, err := newS3Service(config)
	if err != nil {
		return nil, err
	}

	allObjectsList, err := listAllObjects(svc, config)
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, obj := range allObjectsList.Contents {
		keys = append(keys, *obj.Key)
	}

	return keys, nil
}
