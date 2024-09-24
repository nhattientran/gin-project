package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net"
	"sync"
	"time"
)

func (s *Server) recoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := recover(); err != nil {
			s.errorLog.PrintError(fmt.Errorf("%s", err), nil)
			s.serverErrorResponse(c, fmt.Errorf("%s", err))
			c.Set("Connection", "close")
		}
		c.Next()
	}
}

func (s *Server) rateLimit() gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, cli := range clients {
				if time.Since(cli.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			s.serverErrorResponse(c, err)
			return
		}
		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(s.config.limiter.rps), s.config.limiter.burst)}
		}
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			s.rateLimitExceededResponse(c)
			return
		}
		mu.Unlock()

		c.Next()
	}
}
