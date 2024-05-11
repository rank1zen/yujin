package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	ip       string
	port     string
	listener net.Listener
}

// NewServer creates a new server listening on the provided address that responds
// to the http.Handler. It starts the listener, but does not start the server. If
// an empty port is given, the server randomly chooses one.
func NewServer(port string) (*Server, error) {
	addr := fmt.Sprintf(":" + port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}

	return &Server{
		ip:       listener.Addr().(*net.TCPAddr).IP.String(),
		port:     strconv.Itoa(listener.Addr().(*net.TCPAddr).Port),
		listener: listener,
	}, nil
}

func (s *Server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	err := srv.Serve(s.listener)
	if err != nil {
		return fmt.Errorf("ok: %w", err)
	}

	return nil
}

func (s *Server) ServeHTTPHandler(ctx context.Context, handler http.Handler) error {
	return s.ServeHTTP(ctx, &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           handler,
	})
}
