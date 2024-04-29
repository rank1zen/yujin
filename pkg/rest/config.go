package rest

import (
	"github.com/KnutZuidema/golio"
	"github.com/rank1zen/yujin/pkg/database"
)

// Env represents environment configuration for a new rest server
type Env struct {
	database    *database.DB
	golioClient *golio.Client
}
