package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/alexferl/air/asset"
)

type LinodeOpts struct {
	Bucket string
	Region string
}

type Linode struct {
	*LinodeOpts
	s3 Storage
}

func NewLinode(ctx context.Context, opts *LinodeOpts) (Storage, error) {
	ep := fmt.Sprintf("https://%s.linodeobjects.com", opts.Region)
	s3Opts := &S3Opts{
		Bucket:       opts.Bucket,
		StorageClass: "STANDARD",
		Region:       opts.Region,
		Endpoint:     ep,
	}
	s3, err := NewS3(ctx, s3Opts)
	if err != nil {
		return nil, err
	}

	return &Linode{
		LinodeOpts: opts,
		s3:         s3,
	}, nil
}

func (l *Linode) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	return l.s3.Get(ctx, path)
}

func (l *Linode) Put(ctx context.Context, a *asset.Asset) error {
	return l.s3.Put(ctx, a)
}
