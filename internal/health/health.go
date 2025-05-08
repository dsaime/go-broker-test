package health

import (
	"context"
	"fmt"
)

type Healthchecking interface {
	Healthcheck(context.Context) error
}

type Health struct {
	Components map[string]Healthchecking
}

func (h Health) Healthcheck(ctx context.Context) error {
	for name, component := range h.Components {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := component.Healthcheck(ctx); err != nil {
			return fmt.Errorf("healthcheck component %s failed: %v", name, err)
		}
	}

	return nil
}
