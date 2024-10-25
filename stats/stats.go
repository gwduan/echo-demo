package stats

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type Stats struct {
	Uptime   time.Time      `json:"uptime"`
	Requests uint64         `json:"requests"`
	Statuses map[string]int `json:"statuses"`
	URLs     map[string]int `json:"urls"`
	mutex    sync.RWMutex
}

func New() *Stats {
	return &Stats{
		Uptime:   time.Now(),
		Statuses: map[string]int{},
		URLs:     map[string]int{},
	}
}

func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.Requests++
		status := strconv.Itoa(c.Response().Status)
		s.Statuses[status]++
		smu := fmt.Sprintf("%s %-6s:%s",
			status,
			c.Request().Method,
			c.Request().URL.String())
		s.URLs[smu]++
		return nil
	}
}

func (s *Stats) Handler(c echo.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return c.JSON(http.StatusOK, s)
}
