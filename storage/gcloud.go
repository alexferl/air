package storage

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"

	"github.com/alexferl/air/asset"
)

type GCloudOpts struct {
	ProjectId    string
	Bucket       string
	StorageClass string
	Location     string
}

type GCloud struct {
	*GCloudOpts
	bucket *storage.BucketHandle
}

func NewGCloud(ctx context.Context, opts *GCloudOpts) (Storage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(opts.Bucket)
	if _, err := bucket.Attrs(ctx); err != nil {
		log.Info().Msgf("Creating bucket %s in project %s", opts.Bucket, opts.ProjectId)
		if err := bucket.Create(ctx, opts.ProjectId, &storage.BucketAttrs{
			StorageClass: opts.StorageClass,
			Location:     opts.Location,
		}); err != nil {
			return nil, err
		}
	}

	return &GCloud{
		GCloudOpts: opts,
		bucket:     bucket,
	}, nil
}

func (gc *GCloud) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	r, err := gc.bucket.Object(path).NewReader(ctx)
	if err != nil {
		log.Error().Msgf("Failed to read object from bucket: %v", err)
		return nil, err
	}

	return r, nil
}

func (gc *GCloud) Put(ctx context.Context, a *asset.Asset) error {
	wc := gc.bucket.Object(a.Path).NewWriter(ctx)
	if _, err := io.Copy(wc, a.File); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		log.Error().Msgf("Failed to write object to bucket: %v", err)
		return err
	}

	return nil
}
