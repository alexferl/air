package storage

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/alexferl/air/asset"
)

type FilesystemOpts struct {
	Path string
}

type Filesystem struct {
	*FilesystemOpts
}

func NewFilesystem(opts *FilesystemOpts) (Storage, error) {
	return &Filesystem{
		FilesystemOpts: opts,
	}, nil
}

func (fs *Filesystem) Get(_ context.Context, path string) (io.ReadCloser, error) {
	fullPath := fs.getFullPath(path)
	f, err := os.OpenFile(fullPath, os.O_RDONLY, 0o644)
	if err != nil {
		log.Error().Msgf("Failed to open file for reading: %v", err)
		return nil, err
	}

	return f, nil
}

func (fs *Filesystem) Put(_ context.Context, a *asset.Asset) error {
	folderPath := fs.getFullPath(a.PathPrefix)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		log.Error().Msgf("Failed to create folders at path %s: %v", folderPath, err)
		return err
	}

	fullPath := fs.getFullPath(a.Path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			log.Error().Msgf("Failed to open file for writing: %v", err)
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, a.File)
		if err != nil {
			log.Error().Msgf("Failed to write to file: %v", err)
			return err
		}
	}

	return nil
}

func (fs *Filesystem) getFullPath(name string) string {
	fullPath := fs.Path
	if len(fullPath) > 0 && string(fullPath[len(fullPath)-1]) != "/" {
		fullPath += "/"
	}

	fullPath += name

	return fullPath
}
