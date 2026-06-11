package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"yoyo-server/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageProvider interface {
	Save(ctx context.Context, key string, body io.Reader, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
	Name() string
}

type LocalStorageProvider struct {
	uploadDir string
}

func NewLocalStorageProvider(uploadDir string) *LocalStorageProvider {
	return &LocalStorageProvider{uploadDir: uploadDir}
}

func (p *LocalStorageProvider) Save(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	fullPath := filepath.Join(p.uploadDir, key)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	target, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer target.Close()
	if _, err := io.Copy(target, body); err != nil {
		return "", err
	}
	return "/uploads/" + key, nil
}

func (p *LocalStorageProvider) Delete(ctx context.Context, key string) error {
	return os.Remove(filepath.Join(p.uploadDir, key))
}

func (p *LocalStorageProvider) Name() string { return "local" }

type R2StorageProvider struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func NewR2StorageProvider(cfg *config.Config) (*R2StorageProvider, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID),
		}, nil
	})

	awsCfg, err := s3config.LoadDefaultConfig(context.TODO(),
		s3config.WithEndpointResolverWithOptions(r2Resolver),
		s3config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.R2AccessKeyID, cfg.R2SecretAccessKey, "")),
		s3config.WithRegion(cfg.R2Region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)
	return &R2StorageProvider{
		client:    client,
		bucket:    cfg.R2Bucket,
		publicURL: cfg.R2PublicBaseURL,
	}, nil
}

func (p *R2StorageProvider) Save(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	_, err := p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", p.publicURL, key), nil
}

func (p *R2StorageProvider) Delete(ctx context.Context, key string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (p *R2StorageProvider) Name() string { return "r2" }
