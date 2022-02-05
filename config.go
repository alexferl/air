package main

import (
	"time"

	xconfig "github.com/alexferl/golib/config"
	xhttp "github.com/alexferl/golib/http/config"
	xlog "github.com/alexferl/golib/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	Config              *xconfig.Config
	Http                *xhttp.Config
	Logging             *xlog.Config
	MaxFileSize         int64
	FileDownloadTimeout time.Duration
	FileUploadTimeout   time.Duration
	Vips                *Vips
	Storage             *Storage
}

// Vips holds vips specific configuration
type Vips struct {
	ConcurrencyLevel int
	MaxCacheMem      int
	MaxCacheSize     int
	MaxCacheFiles    int
}

type Storage struct {
	Type       string
	Filesystem *Filesystem
	GCloud     *GCloud
	Linode     *Linode
	S3         *S3
}

type Filesystem struct {
	Path string
}

type GCloud struct {
	ProjectId    string
	Bucket       string
	StorageClass string
	Location     string
}

type Linode struct {
	Bucket string
	Region string
}

type S3 struct {
	Bucket       string
	StorageClass string
	Region       string
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	return &Config{
		Config:              xconfig.New(),
		Http:                xhttp.DefaultConfig,
		Logging:             xlog.DefaultConfig,
		MaxFileSize:         10,
		FileUploadTimeout:   time.Second * 60,
		FileDownloadTimeout: time.Second * 60,
		Vips: &Vips{
			ConcurrencyLevel: 1,
			MaxCacheMem:      100 * 1024 * 1024, // 100MB
			MaxCacheSize:     500,
			MaxCacheFiles:    100,
		},
		Storage: &Storage{
			Type: "filesystem",
			Filesystem: &Filesystem{
				Path: "/tmp",
			},
			GCloud: &GCloud{
				ProjectId:    "",
				Bucket:       "",
				Location:     "US",
				StorageClass: "STANDARD",
			},
			Linode: &Linode{
				Bucket: "",
				Region: "",
			},
			S3: &S3{
				Bucket:       "",
				Region:       "",
				StorageClass: "STANDARD",
			},
		},
	}
}

// addFlags adds all the flags from the command line
func (c *Config) addFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&c.FileDownloadTimeout, "file-download-timeout", c.FileDownloadTimeout,
		"File download timeout")
	fs.DurationVar(&c.FileUploadTimeout, "file-upload-timeout", c.FileUploadTimeout, "File upload timeout")
	fs.Int64Var(&c.MaxFileSize, "max-file-size", c.MaxFileSize, "Max file size for uploads in megabytes")

	// Vips
	fs.IntVar(&c.Vips.ConcurrencyLevel, "vips-concurrency-level", c.Vips.ConcurrencyLevel,
		"vips concurrency level")
	fs.IntVar(&c.Vips.MaxCacheMem, "vips-max-cache-mem", c.Vips.MaxCacheMem,
		"Maximum amount of memory in bytes vips is allowed to use for caching")
	fs.IntVar(&c.Vips.MaxCacheSize, "vips-max-cache-size", c.Vips.MaxCacheSize,
		"Maximum amount of objects vips is allowed to keep in cache")
	fs.IntVar(&c.Vips.MaxCacheFiles, "vips-max-cache-files", c.Vips.MaxCacheFiles,
		"Maximum amount of files vips is allowed to keep in cache")

	// Storage
	fs.StringVar(&c.Storage.Type, "storage-type", c.Storage.Type,
		"Storage type to use for assets (filesystem, gcloud, linode, s3)")

	// Filesystem
	fs.StringVar(&c.Storage.Filesystem.Path, "filesystem-path", c.Storage.Filesystem.Path,
		"Filesystem storage path")

	// GCloud
	fs.StringVar(&c.Storage.GCloud.ProjectId, "gcloud-project-id", c.Storage.GCloud.ProjectId,
		"Google Cloud project id")
	fs.StringVar(&c.Storage.GCloud.Bucket, "gcloud-bucket", c.Storage.GCloud.Bucket,
		"Google Cloud storage bucket")
	fs.StringVar(&c.Storage.GCloud.StorageClass, "gcloud-storage-class", c.Storage.GCloud.StorageClass,
		"Google Cloud bucket storage class")
	fs.StringVar(&c.Storage.GCloud.Location, "gcloud-location", c.Storage.GCloud.Location,
		"Google Cloud bucket location")

	// Linode
	fs.StringVar(&c.Storage.Linode.Bucket, "linode-bucket", c.Storage.Linode.Bucket,
		"Linode storage bucket")
	fs.StringVar(&c.Storage.Linode.Region, "linode-region", c.Storage.Linode.Region, "Linode region")

	// S3
	fs.StringVar(&c.Storage.S3.Bucket, "s3-bucket", c.Storage.S3.Bucket, "AWS S3 storage bucket")
	fs.StringVar(&c.Storage.S3.Region, "s3-region", c.Storage.S3.Region, "AWS S3 region")
	fs.StringVar(&c.Storage.S3.StorageClass, "s3-storage-class", c.Storage.S3.StorageClass,
		"AWS S3 storage class")
}

func (c *Config) BindFlags() {
	c.addFlags(pflag.CommandLine)
	c.Logging.BindFlags(pflag.CommandLine)
	c.Http.BindFlags(pflag.CommandLine)

	err := c.Config.BindFlags()
	if err != nil {
		panic(err)
	}

	err = xlog.New(&xlog.Config{
		LogLevel:  viper.GetString("log-level"),
		LogOutput: viper.GetString("log-output"),
		LogWriter: viper.GetString("log-writer"),
	})
	if err != nil {
		panic(err)
	}
}
