package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"

	"github.com/alexferl/air/asset"
)

type S3Opts struct {
	Bucket       string
	StorageClass string
	Region       string
	Endpoint     string
}

type S3 struct {
	*S3Opts
	client *s3.Client
}

func NewS3(ctx context.Context, opts *S3Opts) (Storage, error) {
	var ep string
	if opts.Endpoint == "" {
		ep = fmt.Sprintf("https://s3.%s.amazonaws.com", opts.Region)
	} else {
		ep = opts.Endpoint
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(opts.Region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           ep,
					SigningRegion: opts.Region,
				}, nil
			})),
	)
	if err != nil {
		return nil, err
	}

	var cbc *types.CreateBucketConfiguration
	if strings.Contains(ep, "amazonaws.com") {
		cbc = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(opts.Region),
		}
	}

	client := s3.NewFromConfig(cfg)

	if _, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(opts.Bucket),
	}); err != nil {
		if _, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket:                    aws.String(opts.Bucket),
			CreateBucketConfiguration: cbc,
		}); err != nil {
			return nil, err
		}
	}

	return &S3{
		S3Opts: opts,
		client: client,
	}, nil
}

func (s *S3) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Error().Msgf("Failed to get object from bucket: %v", err)
		return nil, err
	}

	return resp.Body, nil
}

func (s *S3) Put(ctx context.Context, a *asset.Asset) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(s.Bucket),
		Key:          aws.String(a.Path),
		StorageClass: types.StorageClass(s.StorageClass),
		Body:         a.File,
	})
	if err != nil {
		log.Error().Msgf("Failed to put object to bucket: %v", err)
		return err
	}

	return nil
}
