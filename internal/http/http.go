package http

import (
	"context"
	"encoding/json"
	"github.com/adscompass/ip2location/internal/ip2location"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var fn ip2location.ParseFunc

func Listen(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, address string, f ip2location.ParseFunc) {
	fn = f

	h := &handler{}

	server := &http.Server{
		Addr:    address,
		Handler: h,
	}

	go func() {
		log.Printf("http listener %s", address)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("error listen and serve http listener, %v", err)
			ctxCancel()
		}

	}()

	go func() {
		<-ctx.Done()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("error shutdown server, %v", err)
		}
		wg.Done()
	}()
}

type handler struct{}

// IP можно передать:
// - в аргументе запроса ip. Например: domain.com?ip=4.0.0.0
// - в заголовке запроса X-IP2LOCATION-IP
// - в теле POST запроса
func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ipStr := req.URL.Query().Get("ip")
	if ipStr == "" {
		ipStr = req.Header.Get("X-IP2LOCATION-IP")
	}
	if ipStr == "" && req.Method == http.MethodPost {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("error read body, %v", err)
			rw.WriteHeader(500)
			return
		}
		ipStr = string(data)
	}

	if ipStr == "" {
		log.Printf("ip not provided")
		rw.WriteHeader(400)
		return
	}

	rec, err := fn(ipStr)
	if err != nil {
		log.Printf("error parse ip, %v", err)
		rw.WriteHeader(400)
		return
	}

	resp, err := json.Marshal(rec)
	if err != nil {
		log.Printf("error marshal json, %v", err)
		rw.WriteHeader(500)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	_, err = rw.Write(resp)
	if err != nil {
		log.Printf("error write response, %v", err)
		rw.WriteHeader(500)
		return
	}
}
