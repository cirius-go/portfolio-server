package awsutil

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/go-playground/validator/v10"
)

// S3 is a service to interact with AWS S3.
type S3 struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	validator     *validator.Validate
}

// NewS3 creates a new S3 service.
func NewS3(cfg aws.Config) *S3 {
	client := s3.NewFromConfig(cfg)
	return &S3{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		validator:     validator.New(),
	}
}

// Simplifies the one in aws sdk to our needs
type PutObjectInput struct {
	Bucket        string `validate:"required"`
	Key           string `validate:"required"`
	ContentLength int64  `validate:"required"`
	ContentMD5    string
	ContentType   string        `validate:"required"`
	Expires       time.Duration `validate:"required"`
}

// PresignPutObject presigns a put object request
func (s *S3) PresignPutObject(ctx context.Context, params PutObjectInput) (*v4.PresignedHTTPRequest, error) {
	if err := s.validator.Struct(params); err != nil {
		return nil, err
	}

	input := &s3.PutObjectInput{
		Bucket:        &params.Bucket,
		Key:           &params.Key,
		ContentType:   &params.ContentType,
		ContentLength: &params.ContentLength,
		ACL:           types.ObjectCannedACLPublicRead,
	}

	if params.ContentMD5 != "" {
		input.ContentMD5 = &params.ContentMD5
	}

	res, err := s.presignClient.PresignPutObject(ctx, input, s3.WithPresignExpires(params.Expires))
	return res, err
}

// ObjectExists checks if the object exists in the bucket
func (s *S3) ObjectExists(ctx context.Context, bucket, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err == nil {
		return true, nil
	}

	var apierr *smithy.OperationError
	if errors.As(err, &apierr) && strings.Contains(apierr.Err.Error(), "Not Found") {
		return false, nil
	}

	return false, err
}

// MoveObject copies the object from/to given paths and delete it from source
func (s *S3) MoveObject(ctx context.Context, bucket, frompath, topath string) error {
	err := s.CopyObject(ctx, bucket, frompath, topath)
	if err == nil {
		err = s.DeleteObject(ctx, bucket, frompath)
	}
	return err
}

// copy from the same bucket.
func (s *S3) CopyObject(ctx context.Context, bucket, frompath, topath string) error {
	_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(fmt.Sprintf("%v/%v", bucket, frompath)),
		Key:        aws.String(topath),
	})
	return err
}

// DeleteObject deletes the object from the bucket
func (s *S3) DeleteObject(ctx context.Context, bucket, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

// PutObject uploads the object to the bucket
func (s *S3) PutObject(ctx context.Context, bucket, key string, body io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}

// PresignGetObject presigns a get object request
func (s *S3) PresignGetObject(ctx context.Context, bucket, key string, exp time.Duration) (*v4.PresignedHTTPRequest, error) {
	return s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(exp))
}

// ZipChildInput is the input for ZipObjects
type ZipChildInput struct {
	Bucket string
	Key    string
}

// ZipObjects zips the objects and uploads the zip to the bucket
func (s *S3) ZipObjects(ctx context.Context, bucket, key string, children ...*ZipChildInput) (*manager.UploadOutput, error) {
	downloader := manager.NewDownloader(s.client)
	downloader.Concurrency = 2

	archive, err := os.Create(fmt.Sprintf("%s.zip", key))
	if err != nil {
		return nil, err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	for _, child := range children {
		err := zipToFile(ctx, downloader, child.Bucket, child.Key, zipWriter)
		if err != nil {
			return nil, err
		}
	}

	uploader := manager.NewUploader(s.client)
	uploader.Concurrency = 2
	uploadOutput, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   archive,
	})

	if err != nil {
		return nil, err
	}

	return uploadOutput, nil
}

func zipToFile(ctx context.Context, downloader *manager.Downloader, bucket, key string, zipWriter *zip.Writer) error {
	writer, err := zipWriter.Create(key)
	if err != nil {
		return err
	}

	buf := manager.NewWriteAtBuffer([]byte{})
	_, err = downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	_, err = writer.Write(buf.Bytes())
	return err
}
