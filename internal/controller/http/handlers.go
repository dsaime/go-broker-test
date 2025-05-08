package http

import (
	"encoding/json"

	"gitlab.com/digineat/go-broker-test/internal/service"
)

func (c *Controller) Ping(context Context) (any, error) {
	return "pong", nil
}

func (c *Controller) EnqueueTrade(context Context) (any, error) {
	var input service.EnqueueTradeInput
	if err := json.NewDecoder(context.request.Body).Decode(&input); err != nil {
		return nil, err
	}

	if err := c.trades.EnqueueTrade(input); err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Controller) AccountStats(context Context) (any, error) {
	input := service.AccountStatisticsInput{
		Account: context.request.PathValue("acc"),
	}

	out, err := c.trades.AccountStatistics(input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Controller) Healthz(context Context) (any, error) {
	if err := c.healthz.Healthcheck(context.request.Context()); err != nil {
		return nil, err
	}

	return nil, nil
}
