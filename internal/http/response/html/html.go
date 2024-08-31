package html

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/rank1zen/yujin/internal/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func OK(w http.ResponseWriter, r *http.Request, c templ.Component) {
	ctx := r.Context()

	w.Header().Set("ContentType", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	c.Render(ctx, w)
}

func ServerError(w http.ResponseWriter, r *http.Request, c templ.Component, err error) {
	ctx := r.Context()

	logger := logging.FromContext(ctx)

	fields := []zapcore.Field{
		zap.Error(err),
	}

	logger.Error(http.StatusText(http.StatusInternalServerError), fields...)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("ContentType", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate, no-store")
	c.Render(ctx, w)
}

func BadRequest(w http.ResponseWriter, r *http.Request, c templ.Component, err error) {
	ctx := r.Context()

	logger := logging.FromContext(ctx)

	fields := []zapcore.Field{
		zap.Error(err),
	}

	logger.Error(http.StatusText(http.StatusBadRequest), fields...)

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("ContentType", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate, no-store")
	c.Render(ctx, w)
}
