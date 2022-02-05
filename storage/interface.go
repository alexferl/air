package storage

import (
	"context"
	"io"

	"github.com/alexferl/air/asset"
)

type Storage interface {
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	Put(ctx context.Context, a *asset.Asset) error
}
