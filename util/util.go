package util

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

func CreateTempFile() (*os.File, error) {
	f, err := ioutil.TempFile("", "air-")
	if err != nil {
		log.Error().Msgf("Failed to create temp file: %v", err)
		return nil, err
	}
	log.Debug().Msgf("Created temp file %s", f.Name())
	return f, nil
}

func CleanupTempFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Error().Msgf("Failed to close temp file %s: %v", f.Name(), err)
	}
	log.Debug().Msgf("Closed temp file %s", f.Name())

	err = os.Remove(f.Name())
	if err != nil {
		log.Error().Msgf("Failed to remove temp file %s: %v", f.Name(), err)
	}
	log.Debug().Msgf("Removed temp file %s", f.Name())
}

func GetFullPathFromSha256(hash string) (string, error) {
	if len(hash) != sha256.BlockSize {
		return "", errors.New("invalid sha256 checksum")
	}
	return fmt.Sprintf("%s/%s/%s/%s", hash[0:2], hash[2:4], hash[4:6], hash), nil
}
