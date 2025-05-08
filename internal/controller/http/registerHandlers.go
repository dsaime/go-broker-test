package http

func (c *Controller) registerHandlers() {
	c.HandleFunc("/ping", c.Ping)

	c.HandleFunc("POST /trades", c.EnqueueTrade)
	c.HandleFunc("GET /stats/{acc}", c.AccountStats)
	c.HandleFunc("GET /healthz", c.Healthz)
}
