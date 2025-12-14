package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	register(mux, "/api/v1/health", os.Getenv("APP_UPSTREAM_AUTH_URL"), cfg, false)
	register(mux, "/api/v1/auth/", os.Getenv("APP_UPSTREAM_AUTH_URL"), cfg, false)
	register(mux, "/api/v1/schools/", os.Getenv("APP_UPSTREAM_ACADEMIC_URL"), cfg, true)
	register(mux, "/api/v1/classes/", os.Getenv("APP_UPSTREAM_ACADEMIC_URL"), cfg, true)
	register(mux, "/api/v1/subjects/", os.Getenv("APP_UPSTREAM_ACADEMIC_URL"), cfg, true)
	register(mux, "/api/v1/attendance/", os.Getenv("APP_UPSTREAM_ATTENDANCE_URL"), cfg, true)
	register(mux, "/api/v1/grades/", os.Getenv("APP_UPSTREAM_ASSESSMENT_URL"), cfg, true)
	register(mux, "/api/v1/reports/", os.Getenv("APP_UPSTREAM_ASSESSMENT_URL"), cfg, true)
	register(mux, "/api/v1/admissions/", os.Getenv("APP_UPSTREAM_ADMISSION_URL"), cfg, true)
	register(mux, "/api/v1/finance/", os.Getenv("APP_UPSTREAM_FINANCE_URL"), cfg, true)
	register(mux, "/api/v1/notifications/", os.Getenv("APP_UPSTREAM_NOTIFICATION_URL"), cfg, true)
	register(mux, "/api/v1/files/", os.Getenv("APP_UPSTREAM_FILE_URL"), cfg, true)
}

func register(mux *http.ServeMux, prefix, upstream string, cfg config.Config, requireAuth bool) {
	if upstream == "" {
		return
	}
	u, err := url.Parse(upstream)
	if err != nil {
		log.Printf("skip route %s: invalid upstream: %v", prefix, err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	var h http.Handler = withBreaker(fmt.Sprintf("breaker:%s", prefix), proxy)
	if requireAuth {
		h = middleware.AuthWith(cfg.JWTAccessSecret, cfg.JWTIssuer, cfg.JWTAudience, h)
	}
	mux.Handle(prefix, h)
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
			"auth":         os.Getenv("APP_UPSTREAM_AUTH_URL"),
			"academic":     os.Getenv("APP_UPSTREAM_ACADEMIC_URL"),
			"attendance":   os.Getenv("APP_UPSTREAM_ATTENDANCE_URL"),
			"assessment":   os.Getenv("APP_UPSTREAM_ASSESSMENT_URL"),
			"admission":    os.Getenv("APP_UPSTREAM_ADMISSION_URL"),
			"finance":      os.Getenv("APP_UPSTREAM_FINANCE_URL"),
			"notification": os.Getenv("APP_UPSTREAM_NOTIFICATION_URL"),
			"file":         os.Getenv("APP_UPSTREAM_FILE_URL"),
		}
		out := resp{Services: map[string]svcStatus{}}
		for name, base := range services {
			if base == "" {
				out.Services[name] = svcStatus{Up: false, Status: 0, Error: "no_upstream"}
				continue
			}
			target := base
			if !strings.HasSuffix(target, "/") {
				target += "/"
			}
			target += "api/v1/health"
			req, _ := http.NewRequest(http.MethodGet, target, nil)
			res, err := client.Do(req)
			if err != nil {
				out.Services[name] = svcStatus{Up: false, Status: 0, Error: err.Error()}
				continue
			}
			_ = res.Body.Close()
			up := res.StatusCode == http.StatusOK
			out.Services[name] = svcStatus{Up: up, Status: res.StatusCode, Error: ""}
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
