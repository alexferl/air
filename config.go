package main

import (
	xconfig "github.com/alexferl/golib/config"
	xhttp "github.com/alexferl/golib/http/config"
	xlog "github.com/alexferl/golib/log"
	"github.com/spf13/pflag"
)

// Config holds all configuration for our program
type Config struct {
	Config      *xconfig.Config
	Http        *xhttp.Config
	Logging     *xlog.Config
	MaxFileSize int64
	Vips        Vips
	Storage     Storage
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
	Filesystem Filesystem
}

type Filesystem struct {
	Path string
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	c := Config{
		Config:      xconfig.New(),
		Http:        xhttp.DefaultConfig,
		Logging:     xlog.DefaultConfig,
		MaxFileSize: 32,
		Vips: Vips{
			ConcurrencyLevel: 1,
			MaxCacheMem:      100 * 1024 * 1024, // 100MB
			MaxCacheSize:     500,
			MaxCacheFiles:    100,
		},
		Storage: Storage{
			Type: "filesystem",
			Filesystem: Filesystem{
				Path: "/tmp",
			},
		},
	}
	return &c
}

// addFlags adds all the flags from the command line
func (c *Config) addFlags(fs *pflag.FlagSet) {
	fs.Int64Var(&c.MaxFileSize, "max-file-size", c.MaxFileSize,
		"Max file size for uploads in megabytes")
	fs.IntVar(&c.Vips.ConcurrencyLevel, "vips-concurrency-level", c.Vips.ConcurrencyLevel,
		"vips concurrency level")
	fs.IntVar(&c.Vips.MaxCacheMem, "vips-max-cache-mem", c.Vips.MaxCacheMem,
		"Maximum amount of memory in bytes vips is allowed to use for caching")
	fs.IntVar(&c.Vips.MaxCacheSize, "vips-max-cache-size", c.Vips.MaxCacheSize,
		"Maximum amount of objects vips is allowed to keep in cache")
	fs.IntVar(&c.Vips.MaxCacheFiles, "vips-max-cache-files", c.Vips.MaxCacheFiles,
		"Maximum amount of files vips is allowed to keep in cache")
	fs.StringVar(&c.Storage.Type, "storage-type", c.Storage.Type,
		"Storage type to use for assets (filesystem, gcloud)")
	fs.StringVar(&c.Storage.Filesystem.Path, "filesystem-path", c.Storage.Filesystem.Path,
		"Filesystem storage path")
}

func (c *Config) BindFlags() {
	c.addFlags(pflag.CommandLine)
	c.Logging.BindFlags(pflag.CommandLine)
	c.Http.BindFlags(pflag.CommandLine)

	err := c.Config.BindFlags()
	if err != nil {
		panic(err)
	}
}
