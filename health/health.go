package health

import (
	"sync"
	"time"
)

type HealthCheck struct {
	LastSuccess time.Time
	Errors      int
	sync.RWMutex
}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

func (h *HealthCheck) RecordSuccess() {
	h.Lock()
	defer h.Unlock()
	h.LastSuccess = time.Now()
	h.Errors = 0
}

func (h *HealthCheck) RecordError() int {
	h.Lock()
	defer h.Unlock()
	h.Errors++
	return h.Errors
}

func (h *HealthCheck) GetStatus() (time.Time, int) {
	h.RLock()
	defer h.RUnlock()
	return h.LastSuccess, h.Errors
}
