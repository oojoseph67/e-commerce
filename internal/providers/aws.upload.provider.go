package providers

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	appconfig "github.com/oojoseph67/ecommerce/internal/config"
	"github.com/rs/zerolog"
)

type S3Provider struct {
	client   *s3.Client
	tmClient *transfermanager.Client
	bucket   string
	endpoint string
	logger   zerolog.Logger
}

func NewS3Provider(cfg *appconfig.AWSConfig, logger zerolog.Logger) *S3Provider {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyId, cfg.SecretAccessKey, "")))

	if err != nil {
		panic("failed to create aws config " + err.Error())
	}

	// config for localstack
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
			o.UsePathStyle = true
		}
	})

	// new transfer manager client (replaces deprecated manager.Uploader)
	tmClient := transfermanager.New(client)

	return &S3Provider{
		client:   client,
		tmClient: tmClient,
		bucket:   cfg.S3Bucket,
		endpoint: cfg.S3Endpoint,
		logger:   logger,
	}
}

func (p *S3Provider) UploadFile(file *multipart.FileHeader, path string) (url, altText string, err error) {
	src, err := file.Open()
	if err != nil {
		p.logger.Err(err).Str("file", file.Filename).Msg("failed to open file for S3 upload")
		return "", "", err
	}
	defer func() {
		if cerr := src.Close(); cerr != nil {
			p.logger.Warn().Err(cerr).Str("file", file.Filename).Msg("failed to close file source")
		}
	}()

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	path = filepath.Join(filepath.Dir(path), newFileName)

	ctx := context.Background()

	if err := p.ensureBucket(ctx); err != nil {
		p.logger.Err(err).Str("bucket", p.bucket).Msg("failed to create S3 bucket")
		return "", "", err
	}

	_, err = p.tmClient.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
		Body:   src,
	})
	if err != nil {
		p.logger.Err(err).Str("bucket", p.bucket).Str("key", path).Msg("failed to upload to S3")
		return "", "", err
	}

	var urlStr string
	if p.endpoint != "" {
		urlStr = p.endpoint + "/" + p.bucket + "/" + path
	} else {
		urlStr = "https://" + p.bucket + ".s3.amazonaws.com/" + path
	}

	name := fmt.Sprintf("%s%d", file.Filename, time.Now().UnixMilli())

	return urlStr, name, nil
}

func (p *S3Provider) DeleteFile(path string) error {
	ctx := context.Background()
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		p.logger.Err(err).Str("bucket", p.bucket).Str("key", path).Msg("failed to delete from S3")
	}
	return err
}

func (p *S3Provider) ensureBucket(ctx context.Context) error {
	_, err := p.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(p.bucket),
	})
	if err == nil {
		return nil // already exists
	}

	_, err = p.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(p.bucket),
	})

	// Ignore if bucket already exists or is already owned
	var bne *types.BucketAlreadyExists
	var bno *types.BucketAlreadyOwnedByYou
	if errors.As(err, &bne) || errors.As(err, &bno) {
		return nil
	}

	return err
}
