package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	mscfg "memesearch/internal/config"
	"memesearch/internal/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type YaClientS3 struct {
	client *s3.Client
	bucket string
}

func GetClient(ctx context.Context, cfg mscfg.S3Config) (*YaClientS3, error) {
	s3cfg, err := config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
		o.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.Key, cfg.Secret, ""))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("can't create S3 config: %w", err)
	}

	client := s3.NewFromConfig(s3cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://storage.yandexcloud.net")
		o.UsePathStyle = true
		o.Region = "ru-central1"
	})

	return &YaClientS3{
		client: client,
		bucket: cfg.Bucket,
	}, nil

}

func (ya *YaClientS3) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := ya.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(ya.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		var noSuchKeyErr *types.NoSuchKey
		if errors.As(err, &noSuchKeyErr) {
			return nil, models.ErrMediaNotFound
		}
		return nil, fmt.Errorf("can't get S3 object: %w", err)
	}
	return result.Body, nil
}

func (ya *YaClientS3) GetObjectLink(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(ya.client)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     aws.String(ya.bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String("inline"),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return req.URL, nil
}

func (ya *YaClientS3) PutObject(ctx context.Context, name string, body io.Reader) error {
	_, err := ya.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(ya.bucket),
		Key:         aws.String(name),
		Body:        body,
		ContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		return fmt.Errorf("can't put object to S3: %w", err)
	}
	return nil
}
