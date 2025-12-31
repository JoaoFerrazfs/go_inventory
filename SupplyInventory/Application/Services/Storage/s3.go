package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	client  *minio.Client
	bucket  string
	region  string
	baseURL string
}

func NewS3Storage(endpoint, accessKey, secretKey, bucket, region string, useSSL bool) (*S3Storage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	// ensure bucket exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region}); err != nil {
			return nil, err
		}
	}

	baseURL := fmt.Sprintf("%s/%s/", endpoint, bucket)
	s := &S3Storage{client: minioClient, bucket: bucket, region: region, baseURL: baseURL}

	// If configured, set bucket policy to public-read for all objects
	if os.Getenv("MINIO_PUBLIC_READ") == "true" {
		policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucket)
		// try to set policy; ignore error if server doesn't support it
		_ = minioClient.SetBucketPolicy(context.Background(), bucket, policy)
	}

	return s, nil
}

func (s *S3Storage) Upload(path string, content io.Reader) (string, error) {
	ctx := context.Background()
	// no size known; use PutObject with -1 size and UnknownLength
	_, err := s.client.PutObject(ctx, s.bucket, path, content, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	// return the object path
	return path, nil
}

func (s *S3Storage) GetURL(path string) (string, error) {
	// generate presigned URL
	ctx := context.Background()
	reqParams := make(url.Values)
	presigned, err := s.client.PresignedGetObject(ctx, s.bucket, path, time.Minute*60, reqParams)
	if err != nil {
		return "", err
	}
	// If configured for public read, return direct public URL (no presigned)
	if os.Getenv("MINIO_PUBLIC_READ") == "true" {
		public := os.Getenv("MINIO_PUBLIC_ENDPOINT")
		if public != "" {
			// ensure scheme present
			parsedPublic, err := url.Parse(public)
			if err != nil || parsedPublic.Scheme == "" {
				public = "http://" + public
			}
			// return public endpoint + /{bucket}/{path}
			return fmt.Sprintf("%s/%s/%s", strings.TrimRight(public, "/"), s.bucket, strings.TrimLeft(path, "/")), nil
		}
	}

	// fallback to presigned URL
	return presigned.String(), nil
}

func (s *S3Storage) Delete(path string) error {
	ctx := context.Background()
	return s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
}
