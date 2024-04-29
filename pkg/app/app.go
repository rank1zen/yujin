package app

import "github.com/rank1zen/yujin/pkg/database"

type App struct {
        ran bool
}

type AppBuilder struct {
        defaultServerPort uint32
        database database.DB
}

func NewBuilder() *AppBuilder {
        var builder *AppBuilder
        return builder
}

func (a *AppBuilder) Build() (*App) {
        return nil
}

func (a *AppBuilder) WithDefaultPort() *AppBuilder {
        a.defaultServerPort = 8000
        return a
}
