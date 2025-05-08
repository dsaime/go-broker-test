package http

import (
	"net/http"
	"slices"

	"gitlab.com/digineat/go-broker-test/internal/health"
	"gitlab.com/digineat/go-broker-test/internal/service"
)

// Context представляет контекст HTTP-запроса
type Context struct {
	request *http.Request
}

type HandlerFunc func(Context) (any, error)

// Controller обрабатывает HTTP-запросы
type Controller struct {
	trades  *service.Trades
	healthz health.Healthchecking
	http.ServeMux
}

func InitController(trades *service.Trades, healthz health.Healthchecking) *Controller {
	c := &Controller{
		trades:   trades,
		healthz:  healthz,
		ServeMux: http.ServeMux{},
	}
	c.registerHandlers()

	return c
}

func (c *Controller) HandleFunc(pattern string, handlerFunc HandlerFunc, middlewares ...middleware) {
	c.ServeMux.HandleFunc(pattern, c.modulation(chain(handlerFunc, middlewares...)))
}

type middleware func(HandlerFunc) HandlerFunc

func chain(h HandlerFunc, middlewares ...middleware) HandlerFunc {
	for _, mw := range slices.Backward(middlewares) {
		h = mw(h)
	}
	return h
}
