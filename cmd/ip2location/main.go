package main

import (
	"context"
	"github.com/adscompass/ip2location/internal/http"
	"github.com/ip2location/ip2location-go"

	"github.com/cristalhq/aconfig"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type config struct {
	Database     string `default:"IP2LOCATION.BIN" usage:"database file"`
	HTTPListener string `default:"0.0.0.0:8001" usage:"http listener"`
	Verbose      bool   `default:"false" usage:"verbose mode"`
}

var (
	version = "undefined"
)

func main() {
	cfg := &config{}

	if err := aconfig.LoaderFor(cfg).WithEnvPrefix("IP2LOCATION").Build().Load(cfg); err != nil {
		log.Printf("error load config, %v", err)
		os.Exit(1)
	}

	log.Printf("ip2location version %s, config %#v", version, cfg)

	db, err := ip2location.OpenDB(cfg.Database)
	if err != nil {
		log.Printf("error open database, %v", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	http.Listen(ctx, ctxCancel, wg, cfg.HTTPListener, cfg.Verbose, db)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGINT)

	select {
	case <-ch:
		ctxCancel()
	case <-ctx.Done():
	}

	close(ch)

	wg.Wait()

	log.Printf("done")
}
