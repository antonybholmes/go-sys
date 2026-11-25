package aws

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type (
	S3Store struct {
		client *s3.Client
		bucket string
	}
)

const (
	ContentTypeHeader    = "Content-Type"
	ContentTypeTextPlain = "text/plain"
	ContentTypeJSON      = "application/json"
)

func NewS3Store(client *s3.Client, bucket string) *S3Store {
	return &S3Store{
		client: client,
		bucket: bucket,
	}
}

func (s *S3Store) Upload(ctx context.Context,
	body io.ReadSeeker,
	key string,
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

	log.Debug().Msgf("Uploaded S3 object: %s", key)

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

// DownloadLocal downloads an S3 object to a local file and returns the destination path.
// If a dir path is given, file is saved with the S3 key's base name. If the path
// includes a file name, that path is used as the destination file path.
func (s *S3Store) DownloadLocal(ctx context.Context, key string, path string) (string, error) {
	// key file name is unique so can strip it

	// If path ends with slash, treat as folder and append the key's base name.
	// Otherwise, treat path as the full destination file path.
	if strings.HasSuffix(path, "/") {
		path = filepath.Join(path, filepath.Base(key))
	}

	// download file to mem
	src, err := s.Download(ctx, key)
	if err != nil {
		return "", err
	}

	defer src.Close()

	// Make sure the destination directory exists.
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, sys.FileModeOwnerFullGroupOtherReadEnter); err != nil {
		return "", err
	}

	// create local file
	localFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	// copy downloaded content to local file
	_, err = io.Copy(localFile, src)
	if err != nil {
		return "", err
	}

	return path, nil
}

// DownloadZip downloads the specified S3 objects into a ZIP archive written to writer.
func (s *S3Store) DownloadZip(ctx context.Context, writer io.Writer, keys []string) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	for _, key := range keys {
		src, err := s.Download(ctx, key)
		if err != nil {
			return err
		}

		// header := &zip.FileHeader{
		// 	Name:   key,
		// 	Method: zip.Deflate,
		// }
		// entry, err := zipWriter.CreateHeader(header)
		// if err != nil {
		// 	src.Close()
		// 	return err
		// }

		// We use path instead of filepath because S3 keys always use forward slashes.
		zipFile, err := zipWriter.Create(path.Base(key))
		if err != nil {
			src.Close()
			return err
		}

		_, copyErr := io.Copy(zipFile, src)
		closeErr := src.Close()

		if copyErr != nil {
			return copyErr
		}

		if closeErr != nil {
			return closeErr
		}
	}

	return nil
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
