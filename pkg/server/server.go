package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/rank1zen/yujin/pkg/logging"
)

type Server struct {
	ip       string
	port     string
	listener net.Listener
}

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
	log := logging.FromContext(ctx).Sugar()

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		log.Debug("contex closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), 5 * time.Second)
		defer done()

		log.Debug("shutting down")
		errCh <-srv.Shutdown(shutdownCtx)
	}()

	// Serve blocks until the provided context is closed
	log.Infof("serving on port: %s", s.port)
	err := srv.Serve(s.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	log.Debug("stoped")

	err = <-errCh 
	if err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}

func (s *Server) ServeHTTPHandler(ctx context.Context, handler http.Handler) error {
	return s.ServeHTTP(ctx, &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           handler,
	})
}
