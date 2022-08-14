package main

import (
	"fmt"
	"log"
	"net/http"
	"test-assignment-cookie-sync/config"
	"test-assignment-cookie-sync/handler"
	"test-assignment-cookie-sync/service"
	"time"

	"github.com/rs/cors"
)

type httpService struct {
	*config.HTTPConfig
	router *http.ServeMux
}

func newHTTPService(cfg *config.HTTPConfig) *httpService {
	return &httpService{
		HTTPConfig: cfg,
		router:     http.NewServeMux(),
	}
}

func (h *httpService) registerRoutes(sc *service.CookieS, ss *service.SyncS) {
	h.registerSync(sc, ss)
}

func (h *httpService) registerSync(sc *service.CookieS, ss *service.SyncS) {
	ch := handler.NewCookieHandler(sc)
	sh := handler.NewCookieSyncHandler(ss)

	h.register(http.MethodGet, "/api/cookie", ch.Cookie)
	h.register(http.MethodPut, "/api/sync", sh.Sync)
}

func (h *httpService) register(method, path string, handler http.HandlerFunc) {
	timeout, _ := time.ParseDuration(h.HTTPConfig.Timeout)
	h.router.Handle(path, http.TimeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("incoming connection from %s to %s", r.RemoteAddr, path)
		handler.ServeHTTP(w, r)
	}), timeout, "request canceled by timeout"))
}

func (h *httpService) run() error {
	var listenAddress = fmt.Sprintf("%s:%d", h.HTTPConfig.Host, h.HTTPConfig.Port)
	log.Printf("service listening on %s", listenAddress)
	return http.ListenAndServe(listenAddress, cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:9003"},
		AllowedHeaders: []string{"*"},
	}).Handler(h.router))
}
