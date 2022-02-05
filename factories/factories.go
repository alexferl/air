package factories

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/alexferl/air/storage"
)

func Storage(ctx context.Context, storageType string) (storage.Storage, error) {
	fsConfig := &storage.FilesystemOpts{Path: viper.GetString("filesystem-path")}

	log.Info().Msgf("Using storage type '%s'", storageType)

	switch storageType {
	case "filesystem":
		return storage.NewFilesystem(fsConfig)
	case "gcloud":
		config := &storage.GCloudOpts{
			ProjectId:    viper.GetString("gcloud-project-id"),
			Bucket:       viper.GetString("gcloud-bucket"),
			StorageClass: viper.GetString("gcloud-storage-class"),
			Location:     viper.GetString("gcloud-location"),
		}
		return storage.NewGCloud(ctx, config)
	case "linode":
		config := &storage.LinodeOpts{
			Bucket: viper.GetString("linode-bucket"),
			Region: viper.GetString("linode-region"),
		}
		return storage.NewLinode(ctx, config)
	case "s3":
		config := &storage.S3Opts{
			Bucket:       viper.GetString("s3-bucket"),
			StorageClass: viper.GetString("s3-storage-class"),
			Region:       viper.GetString("s3-region"),
		}
		return storage.NewS3(ctx, config)
	default:
		log.Warn().Msgf("Unknown storage type '%s'. Falling back to 'filesystem'", storageType)
		return storage.NewFilesystem(fsConfig)
	}
}
