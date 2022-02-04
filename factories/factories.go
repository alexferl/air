package factories

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/alexferl/air/storage"
)

func StorageFactory(storageType string) storage.Storage {
	path := viper.GetString("filesystem-path")

	switch storageType {
	case "filesystem":
		log.Info().Msgf("Using storage type '%s'", storageType)
		return storage.NewFilesystem(path)
	default:
		log.Warn().Msgf("Unknown storage type '%s'. Falling back to 'filesystem'", storageType)
		return storage.NewFilesystem(path)
	}
}
