package main

import (
	"context"
	"net/http"
	"time"

	"github.com/alexferl/golib/http/router"
	"github.com/alexferl/golib/http/server"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/alexferl/air/factories"
	"github.com/alexferl/air/handlers"
)

func main() {
	c := NewConfig()
	c.BindFlags()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	storage, err := factories.Storage(ctx, viper.GetString("storage-type"))
	if err != nil {
		panic(err)
	}

	s := server.New()
	h := &handlers.Handler{Storage: storage}
	r := &router.Router{
		Routes: []router.Route{
			{"Root", http.MethodGet, "/", h.Root},
			{"Asset", http.MethodGet, "/assets/:id", h.Asset},
			{"Stats", http.MethodGet, "/stats", h.Stats},
			{"Upload", http.MethodPost, "/upload", h.Upload},
			{"FavIcon", http.MethodGet, "/favicon.ico", func(c echo.Context) error {
				return c.String(http.StatusOK, "")
			}},
		},
	}

	conf := vips.Config{
		CollectStats:     true,
		ConcurrencyLevel: viper.GetInt("vips-concurrency-level"),
		MaxCacheMem:      viper.GetInt("vips-max-cache-mem"),
		MaxCacheSize:     viper.GetInt("vips-max-cache-size"),
		MaxCacheFiles:    viper.GetInt("vips-max-cache-files"),
	}
	vips.LoggingSettings(nil, vips.LogLevelWarning)
	vips.Startup(&conf)
	defer vips.Shutdown()

	s.Start(r)
}
