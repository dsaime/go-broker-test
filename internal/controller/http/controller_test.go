package http

import (
	"net/http"
	"testing"
)

func Test_Server(t *testing.T) {
	t.Run("это http сервер", func(t *testing.T) {
		var _ http.Handler = new(Controller)
	})
	t.Run("context с доступом к данным", func(t *testing.T) {
		_ = &Context{
			request: (*http.Request)(nil),
		}
	})
	t.Run("есть метод modulation", func(t *testing.T) {
		// Функция для преобразования controller.HandlerFunc в тип http.HandlerFunc
		var _ interface {
			modulation(HandlerFunc) http.HandlerFunc
		} = new(Controller)
	})
	t.Run("наличие методов", func(t *testing.T) {
		// Контроллер содержит методы для обработки следующих запросов:
		var _ interface {
			EnqueueTrade(Context) (any, error)
			AccountStats(Context) (any, error)
			Healthz(Context) (any, error)
		} = new(Controller)
	})
}
