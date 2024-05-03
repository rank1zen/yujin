package rest

import (
	"context"
	"fmt"

	"github.com/KnutZuidema/golio"
	"github.com/rank1zen/yujin/pkg/database"
)

type Server struct {
	cfg *Config
	db    database.DB
	riot *golio.Client
}

// Config represents a configuration for a new rest server
type Config struct {
}

// Env represents environment configuration for a new rest server

// NewServer creates a new rest server from a config
func NewServer(ctx context.Context, cfg *Config) (*Server, error) {
	// cfg is supposed to be some config to the server but we are not implementing it yet

	if env.db == nil {
		return nil, fmt.Errorf("missing database in server environment")
	}

	if env.riot == nil {
		return nil, fmt.Errorf("missing golio client in server environment")
	}

	return &Server{
		cfg: cfg,
	}, nil
}
