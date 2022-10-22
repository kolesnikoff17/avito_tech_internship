package httpserver

import (
	"net"
	"time"
)

// Option is a type of functions-setters
type Option func(*Server)

// Port sets up server port
func Port(port string) Option {
	return func(s *Server) {
		if port != "" {
			s.server.Addr = net.JoinHostPort("", port)
		}
	}
}

// ReadTimeout sets up server ReadTimeout
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		if timeout.Seconds() != 0 {
			s.server.ReadTimeout = timeout
		}
	}
}

// WriteTimeout sets up server WriteTimeout
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		if timeout.Seconds() != 0 {
			s.server.WriteTimeout = timeout
		}
	}
}

// ShutdownTimeout sets up server ShutdownTimeout
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		if timeout.Seconds() != 0 {
			s.shutdownTimeout = timeout
		}
	}
}
