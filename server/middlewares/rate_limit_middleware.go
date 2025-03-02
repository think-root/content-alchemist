package middlewares

import (
	"content-alchemist/config"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type visitor struct {
	lastSeen  time.Time
	requests  int
	resetTime time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func getClientIP(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		return strings.Split(forwardedFor, ",")[0]
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	if ip == "::1" || ip == "127.0.0.1" {
		return "localhost"
	}

	return ip
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		mu.Lock()
		v, exists := visitors[ip]
		now := time.Now()

		if !exists {
			visitors[ip] = &visitor{
				lastSeen:  now,
				requests:  1,
				resetTime: now.Add(time.Minute),
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if now.After(v.resetTime) {
			v.requests = 1
			v.resetTime = now.Add(time.Minute)
		} else if v.requests >= config.RATE_LIMIT {
			mu.Unlock()
			http.Error(w, "Rate limit exceeded. Try again in a minute.", http.StatusTooManyRequests)
			return
		} else {
			v.requests++
		}

		v.lastSeen = now
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
