package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = ":8080"
	defaultShutdownTimeout = 3 * time.Second
)

// Server keeps http.Server and some useful helpers
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// New is a constructor for Server
func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}
	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}
	for _, opt := range opts {
		opt(s)
	}

	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()

	return s
}

// Notify returns server's error chan
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown sets up timer for shutdown and sends a signal through context.Context
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
