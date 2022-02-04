package storage

import (
	"io"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/alexferl/air/asset"
)

type Filesystem struct {
	Path string
}

func NewFilesystem(path string) Storage {
	return &Filesystem{
		Path: path,
	}
}

func (fs *Filesystem) Get(name string) (io.ReadCloser, error) {
	fullPath := fs.getFullPath(name)
	f, err := os.OpenFile(fullPath, os.O_RDONLY, 0644)
	if err != nil {
		log.Error().Msgf("Failed to open file for reading: %v", err)
		return nil, err
	}

	return f, nil
}

func (fs *Filesystem) Put(a *asset.Asset) error {
	folderPath := fs.getFullPath(a.PathPrefix)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		log.Error().Msgf("Failed to create folders at path %s: %v", folderPath, err)
		return err
	}

	fullPath := fs.getFullPath(a.Path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
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
	var fullPath = fs.Path
	if len(fullPath) > 0 && string(fullPath[len(fullPath)-1]) != "/" {
		fullPath += "/"
	}

	fullPath += name

	return fullPath
}
