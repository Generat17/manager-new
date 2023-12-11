package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type handler interface {
	InitRouter() http.Handler
}

type Server struct {
	server *http.Server
}

func New(handler handler, port string) *Server {
	return &Server{
		server: &http.Server{
			Addr:           ":" + port,
			Handler:        handler.InitRouter(),
			MaxHeaderBytes: 1 << 20, // 1 MB
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10000 * time.Second,
		},
	}
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		go func() {
			err := s.server.ListenAndServe()
			if err != nil && err.Error() != "http: Server closed" {
				log.Print("failed to start the default server")
			}
		}()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			log.Print("failed to shutdown the default server")
		}
	}()
}
