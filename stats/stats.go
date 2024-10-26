package stats

import (
	"context"
	"echo-demo/redis"
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

type RedisStats struct {
	Requests int            `json:"requests"`
	Statuses map[string]int `json:"statuses"`
	URLs     map[string]int `json:"urls"`
}

type AllStats struct {
	Local  *Stats      `json:"Local"`
	Global *RedisStats `json:"Global"`
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
			client := redis.Client()
			pipe := client.TxPipeline()

			pipe.Incr(ctx, "Requests")
			pipe.HIncrBy(ctx, "Statuses", status, 1)
			pipe.HIncrBy(ctx, "URLs", smu, 1)

			if _, err := pipe.Exec(ctx); err != nil {
				c.Echo().Logger.Debug(err)
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
	client := redis.Client()

	requests, err := client.Get(ctx, "Requests").Result()
	if err != nil {
		c.Echo().Logger.Debug(err)
		return err
	}
	statuses, err := client.HGetAll(ctx, "Statuses").Result()
	if err != nil {
		c.Echo().Logger.Debug(err)
		return err
	}
	urls, err := client.HGetAll(ctx, "URLs").Result()
	if err != nil {
		c.Echo().Logger.Debug(err)
		return err
	}

	rs := &RedisStats{
		Statuses: map[string]int{},
		URLs:     map[string]int{},
	}
	rs.Requests, _ = strconv.Atoi(requests)
	for k, v := range statuses {
		rs.Statuses[k], _ = strconv.Atoi(v)
	}
	for k, v := range urls {
		rs.URLs[k], _ = strconv.Atoi(v)
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return c.JSON(http.StatusOK, AllStats{Local: s, Global: rs})
}
