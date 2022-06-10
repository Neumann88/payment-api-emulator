package server

import (
	"context"
	"net/http"
	"time"
)

type server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewHttpServer(handler http.Handler, port string, readTimeout, writeTimeout, shutdownTime time.Duration) *server {
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	}

	serv := &server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: shutdownTime * time.Second,
	}

	serv.start()

	return serv
}

func (s *server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *server) Notify() <-chan error {
	return s.notify
}

func (s *server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
