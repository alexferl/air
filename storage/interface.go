package storage

import (
	"io"

	"github.com/alexferl/air/asset"
)

type Storage interface {
	Get(string) (io.ReadCloser, error)
	Put(a *asset.Asset) error
}
