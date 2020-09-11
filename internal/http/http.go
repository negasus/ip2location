package http

import (
	"context"
	"encoding/json"
	"github.com/ip2location/ip2location-go"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var (
	db *ip2location.DB
)

func Listen(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, address string, ip2db *ip2location.DB) {
	db = ip2db

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	server := &http.Server{
		Addr:    address,
		Handler: mux,
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

// IP можно передать:
// - в аргументе запроса ip. Например: domain.com?ip=4.0.0.0
// - в заголовке запроса X-IP2LOCATION-IP
// - в теле POST запроса
func getIP(req *http.Request) string {
	ip := req.URL.Query().Get("ip")
	if ip != "" {
		return ip
	}
	ip = req.Header.Get("X-IP2LOCATION-IP")
	if ip != "" {
		return ip
	}
	if req.Method == http.MethodPost {
		data, err := ioutil.ReadAll(req.Body)
		if err == nil && len(data) > 0 {
			return string(data)
		}
	}

	return ""
}

func handler(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	rec, err := db.Get_all(getIP(req))
	if err != nil {
		log.Printf("error fetch data, %v", err)
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
