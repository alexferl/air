package main

import (
	"net/http"

	"github.com/alexferl/golib/http/router"
	"github.com/alexferl/golib/http/server"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/spf13/viper"

	"github.com/alexferl/air/factories"
	"github.com/alexferl/air/handlers"
)

func main() {
	c := NewConfig()
	c.BindFlags()

	s := server.New()
	r := &router.Router{}
	h := &handlers.Handler{
		Storage: factories.StorageFactory(viper.GetString("storage-type")),
	}
	r.Routes = []router.Route{
		{"Root", http.MethodGet, "/", h.Root},
		{"Asset", http.MethodGet, "/:id", h.Asset},
		{"Stats", http.MethodGet, "/stats", h.Stats},
		{"Upload", http.MethodPost, "/upload", h.Upload},
	}

	conf := vips.Config{
		CollectStats:     true,
		ConcurrencyLevel: viper.GetInt("vips-concurrency-level"),
		MaxCacheMem:      viper.GetInt("vips-max-cache-mem"),
		MaxCacheSize:     viper.GetInt("vips-max-cache-size"),
		MaxCacheFiles:    viper.GetInt("vips-max-cache-files"),
	}
	vips.Startup(&conf)
	defer vips.Shutdown()

	s.Start(r)
}
