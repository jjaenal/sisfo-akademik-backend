package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	httpx "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/middleware"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
)

var perPrefixLimits = map[string]int{
	"/api/v1/auth/": 5,
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	l, err := logger.New(cfg.Env)
	if err != nil {
		panic(err)
	}
	redis := redisutil.New(cfg.RedisAddr)
	limiter := redisutil.NewLimiterFromCounter(redis.Raw())

	mux := http.NewServeMux()
	registerRoutes(mux, cfg)

	h := middleware.Recover(
		middleware.RequestID(
			middleware.Logging(l,
				middleware.RateLimitByPolicy(limiter, 100, 30, perPrefixLimits,
					middleware.CORS(cfg.CORSAllowedOrigins, mux),
				),
			),
		),
	)

	s := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.HTTPPort),
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := s.ListenAndServe(); err != nil {
		l.Error("gateway failed")
		panic(err)
	}
}

func registerRoutes(mux *http.ServeMux, cfg config.Config) {
	mux.Handle("/api/v1/gateway/health", gatewayHealthHandler(cfg))
	registerLBEnv(mux, "/api/v1/health", "APP_UPSTREAM_AUTH", cfg, false)
	registerLBEnv(mux, "/api/v1/auth/", "APP_UPSTREAM_AUTH", cfg, false)
	registerLBEnv(mux, "/api/v1/users/", "APP_UPSTREAM_AUTH", cfg, true)
	registerLBEnv(mux, "/api/v1/users", "APP_UPSTREAM_AUTH", cfg, true)
	registerLBEnv(mux, "/api/v1/schools/", "APP_UPSTREAM_ACADEMIC", cfg, true)
	registerLBEnv(mux, "/api/v1/classes/", "APP_UPSTREAM_ACADEMIC", cfg, true)
	registerLBEnv(mux, "/api/v1/subjects/", "APP_UPSTREAM_ACADEMIC", cfg, true)
	registerLBEnv(mux, "/api/v1/attendance/", "APP_UPSTREAM_ATTENDANCE", cfg, true)
	registerLBEnv(mux, "/api/v1/grades/", "APP_UPSTREAM_ASSESSMENT", cfg, true)
	registerLBEnv(mux, "/api/v1/reports/", "APP_UPSTREAM_ASSESSMENT", cfg, true)
	registerLBEnv(mux, "/api/v1/admissions/", "APP_UPSTREAM_ADMISSION", cfg, true)
	registerLBEnv(mux, "/api/v1/finance/", "APP_UPSTREAM_FINANCE", cfg, true)
	registerLBEnv(mux, "/api/v1/notifications/", "APP_UPSTREAM_NOTIFICATION", cfg, true)
	registerLBEnv(mux, "/api/v1/files/", "APP_UPSTREAM_FILE", cfg, true)
}

func registerLBEnv(mux *http.ServeMux, prefix, envBase string, cfg config.Config, requireAuth bool) {
	upstreams := parseUpstreams(envBase)
	if len(upstreams) == 0 {
		return
	}
	var h http.Handler = newRoundRobinProxy(prefix, upstreams)
	if requireAuth {
		h = middleware.AuthWith(cfg.JWTAccessSecret, cfg.JWTIssuer, cfg.JWTAudience, h)
	}
	mux.Handle(prefix, h)
}

func parseUpstreams(envBase string) []string {
	// Prefer comma-separated list in *_URLS, fallback to single *_URL
	urls := os.Getenv(envBase + "_URLS")
	if urls != "" {
		parts := strings.Split(urls, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
		return out
	}
	u := os.Getenv(envBase + "_URL")
	if u == "" {
		return nil
	}
	return []string{u}
}

type recWriter struct {
	http.ResponseWriter
	status int
}

func (rw *recWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func gatewayHealthHandler(cfg config.Config) http.Handler {
	type svcStatus struct {
		Up     bool   `json:"up"`
		Status int    `json:"status"`
		Error  string `json:"error"`
	}
	type resp struct {
		Services map[string]svcStatus `json:"services"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{Timeout: 3 * time.Second}
		services := map[string]string{
			"auth":         "", // resolved via parseUpstreams below
			"academic":     "",
			"attendance":   "",
			"assessment":   "",
			"admission":    "",
			"finance":      "",
			"notification": "",
			"file":         "",
		}
		out := resp{Services: map[string]svcStatus{}}
		for name := range services {
			var envBase string
			switch name {
			case "auth":
				envBase = "APP_UPSTREAM_AUTH"
			case "academic":
				envBase = "APP_UPSTREAM_ACADEMIC"
			case "attendance":
				envBase = "APP_UPSTREAM_ATTENDANCE"
			case "assessment":
				envBase = "APP_UPSTREAM_ASSESSMENT"
			case "admission":
				envBase = "APP_UPSTREAM_ADMISSION"
			case "finance":
				envBase = "APP_UPSTREAM_FINANCE"
			case "notification":
				envBase = "APP_UPSTREAM_NOTIFICATION"
			case "file":
				envBase = "APP_UPSTREAM_FILE"
			}
			ups := parseUpstreams(envBase)
			if len(ups) == 0 {
				out.Services[name] = svcStatus{Up: false, Status: 0, Error: "no_upstream"}
				continue
			}
			up := false
			lastStatus := 0
			lastErr := ""
			for _, base := range ups {
				target := base
				if !strings.HasSuffix(target, "/") {
					target += "/"
				}
				target += "api/v1/health"
				req, _ := http.NewRequest(http.MethodGet, target, nil)
				res, err := client.Do(req)
				if err != nil {
					lastErr = err.Error()
					continue
				}
				_ = res.Body.Close()
				lastStatus = res.StatusCode
				if res.StatusCode == http.StatusOK {
					up = true
					lastErr = ""
					break
				}
			}
			out.Services[name] = svcStatus{Up: up, Status: lastStatus, Error: lastErr}
		}
		httpx.Success(w, out)
	})
}

type breaker struct {
	failures    int
	openedUntil time.Time
	threshold   int
	openFor     time.Duration
}

func newBreaker(threshold int, openFor time.Duration) *breaker {
	return &breaker{threshold: threshold, openFor: openFor}
}

func (b *breaker) allow() bool {
	return time.Now().After(b.openedUntil)
}

func (b *breaker) recordFailure() {
	b.failures++
	if b.failures >= b.threshold {
		b.openedUntil = time.Now().Add(b.openFor)
		b.failures = 0
	}
}

func withBreaker(name string, next http.Handler) http.Handler {
	br := newBreaker(5, 30*time.Second)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !br.allow() {
			httpx.Error(w, http.StatusServiceUnavailable, "1001", "Service unavailable", nil)
			return
		}
		rr := &recWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rr, r)
		if rr.status >= 500 {
			br.recordFailure()
		}
	})
}

type rrProxy struct {
	proxies    []*httputil.ReverseProxy
	index      uint32
	breakers   []*breaker
	prefixName string
}

func newRoundRobinProxy(prefix string, upstreams []string) http.Handler {
	var proxies []*httputil.ReverseProxy
	var breakers []*breaker
	for _, u := range upstreams {
		parsed, err := url.Parse(u)
		if err != nil {
			log.Printf("skip invalid upstream for %s: %s, err=%v", prefix, u, err)
			continue
		}
		p := httputil.NewSingleHostReverseProxy(parsed)
		// Mark failures when transport errors occur
		b := newBreaker(5, 30*time.Second)
		p.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			b.recordFailure()
			httpx.Error(w, http.StatusBadGateway, "6002", "Upstream error", []map[string]string{
				{"upstream": parsed.String(), "error": e.Error()},
			})
		}
		proxies = append(proxies, p)
		breakers = append(breakers, b)
	}
	if len(proxies) == 0 {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpx.Error(w, http.StatusServiceUnavailable, "1001", "No valid upstreams", nil)
		})
	}
	rr := &rrProxy{proxies: proxies, breakers: breakers, prefixName: prefix}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startIdx := int(atomic.AddUint32(&rr.index, 1)) % len(rr.proxies)
		// Find first allowed (circuit closed)
		idx := startIdx
		for i := 0; i < len(rr.proxies); i++ {
			j := (startIdx + i) % len(rr.proxies)
			if rr.breakers[j].allow() {
				idx = j
				break
			}
		}
		// Capture status to record failures
		rrw := &recWriter{ResponseWriter: w, status: 200}
		rr.proxies[idx].ServeHTTP(rrw, r)
		if rrw.status >= 500 {
			rr.breakers[idx].recordFailure()
		}
	})
}
