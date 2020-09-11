package main

import (
	"context"
	"github.com/adscompass/ip2location/internal/http"
	"github.com/ip2location/ip2location-go"

	//"github.com/adscompass/ip2location/internal/ip2location"
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
}

func main() {
	cfg := &config{}

	if err := aconfig.LoaderFor(cfg).WithEnvPrefix("IP2LOCATION").Build().Load(cfg); err != nil {
		log.Printf("error load config, %v", err)
		os.Exit(1)
	}

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
	http.Listen(ctx, ctxCancel, wg, cfg.HTTPListener, db)

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
