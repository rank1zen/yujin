package rest

import (
	"context"
	"fmt"

	"github.com/KnutZuidema/golio"
	"github.com/labstack/echo/v4"
)

type Server struct {
	config      *Config
	database    *database.Query
	golioClient *golio.Client
}

// Config represents a configuration for a new rest server
type Config struct {
}

// NewServer creates a new rest server from a config
func NewServer(ctx context.Context, cfg *Config, env *Env) (*Server, error) {
	// cfg is supposed to be some config to the server but we are not implementing it yet

	if env.database == nil {
		return nil, fmt.Errorf("missing database in server environment")
	}

	if env.golioClient == nil {
		return nil, fmt.Errorf("missing golio client in server environment")
	}

	return &Server{
		config:      cfg,
		database:    database.NewQuery(env.database),
		golioClient: env.golioClient,
	}, nil
}

func (s *Server) Routes(e *echo.Echo) {
}
