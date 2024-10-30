package stats

import (
	"context"
	"echo-demo/vk"
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

type ValkeyStats struct {
	Requests int64            `json:"requests"`
	Statuses map[string]int64 `json:"statuses"`
	URLs     map[string]int64 `json:"urls"`
}

type AllStats struct {
	Local  *Stats       `json:"Local"`
	Global *ValkeyStats `json:"Global"`
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

		status := strconv.Itoa(c.Response().Status)
		smu := fmt.Sprintf("%s %-6s:%s",
			status,
			c.Request().Method,
			c.Request().URL.String())

		go func() {
			ctx := context.Background()
			client := vk.Client()

			for _, resp := range client.DoMulti(ctx,
				client.B().Incr().Key("Requests").Build(),
				client.B().Hincrby().Key("Statuses").Field(status).Increment(1).Build(),
				client.B().Hincrby().Key("URLs").Field(smu).Increment(1).Build()) {
				if err := resp.Error(); err != nil {
					c.Echo().Logger.Debug(err)
				}
			}
		}()

		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.Requests++
		s.Statuses[status]++
		s.URLs[smu]++

		return nil
	}
}

func (s *Stats) Handler(c echo.Context) error {
	ctx := context.Background()
	client := vk.Client()

	requests, err := client.Do(ctx, client.B().Get().Key("Requests").Build()).AsInt64()
	if err != nil {
		c.Echo().Logger.Debug(err)
		return err
	}
	statuses, err := client.Do(ctx, client.B().Hgetall().Key("Statuses").Build()).AsIntMap()
	if err != nil {
		c.Echo().Logger.Debug(err)
		return err
	}
	urls, err := client.Do(ctx, client.B().Hgetall().Key("URLs").Build()).AsIntMap()
	if err != nil {
		c.Echo().Logger.Debug(err)
		return err
	}

	rs := &ValkeyStats{
		Requests: requests,
		Statuses: statuses,
		URLs:     urls,
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return c.JSON(http.StatusOK, AllStats{Local: s, Global: rs})
}
