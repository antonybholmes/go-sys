package aws

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type (
	S3Store struct {
		client *s3.Client
		bucket string
	}
)

func NewS3Store(client *s3.Client, bucket string) *S3Store {
	return &S3Store{
		client: client,
		bucket: bucket,
	}
}

func (s *S3Store) Upload(ctx context.Context,
	key string,
	body io.ReadSeeker,
	contentType string) error {
	input := s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   body,
	}

	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	_, err := s.client.PutObject(ctx, &input)

	if err != nil {
		return fmt.Errorf("put S3 object: %w", err)
	}

	return nil
}

func (s *S3Store) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	input := s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	output, err := s.client.GetObject(ctx, &input)

	if err != nil {
		return nil, fmt.Errorf("get S3 object: %w", err)
	}

	return output.Body, nil
}

func (s *S3Store) Delete(ctx context.Context, key string) error {
	input := s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, &input)

	if err != nil {
		return fmt.Errorf("delete S3 object: %w", err)
	}

	return nil
}

func (s *S3Store) List(ctx context.Context, prefix string) ([]string, error) {
	input := s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}

	output, err := s.client.ListObjectsV2(ctx, &input)

	if err != nil {
		return nil, fmt.Errorf("list S3 objects: %w", err)
	}

	keys := make([]string, 0, len(output.Contents))
	for _, object := range output.Contents {
		keys = append(keys, *object.Key)
	}

	return keys, nil
}

func (s *S3Store) Copy(ctx context.Context, sourceKey, destinationKey string) error {
	input := s3.CopyObjectInput{
		Bucket:     aws.String(s.bucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", s.bucket, sourceKey)),
		Key:        aws.String(destinationKey),
	}

	_, err := s.client.CopyObject(ctx, &input)

	if err != nil {
		return fmt.Errorf("copy S3 object: %w", err)
	}

	return nil
}

func (s *S3Store) Move(ctx context.Context, sourceKey, destinationKey string) error {
	if err := s.Copy(ctx, sourceKey, destinationKey); err != nil {
		return fmt.Errorf("move S3 object: %w", err)
	}

	if err := s.Delete(ctx, sourceKey); err != nil {
		return fmt.Errorf("move S3 object: %w", err)
	}

	return nil
}

func (s *S3Store) Bucket() string {
	return s.bucket
}
