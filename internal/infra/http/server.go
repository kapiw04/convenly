package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func Start(srv *http.Server) {
	log.Printf("Server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func Stop(ctx context.Context, srv *http.Server) {
	log.Println("Shutting down server...")
	_ = srv.Shutdown(ctx)
}
