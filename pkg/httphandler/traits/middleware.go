package traits

type MiddlewareHandler struct {
	handler []func(interface{}) error
}

func (h *MiddlewareHandler) AddMiddleware(m func(interface{}) error) {
	h.handler = append(h.handler, m)
}

func (h *MiddlewareHandler) Handle() (err error) {
	for _, handler := range h.handler {
		if err = handler(h); err != nil {
			return
		}
	}

	return
}
