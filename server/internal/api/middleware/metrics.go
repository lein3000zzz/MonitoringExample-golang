package middleware

import (
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by method, endpoint, and status",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_latency_seconds",
			Help:    "Latency of HTTP requests in seconds",
			Buckets: []float64{0.005, 0.01, 0.02, 0.04}, // Бакеты для SLA
		},
		[]string{"method", "endpoint"},
	)

	externalServiceLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_service_latency_seconds",
			Help:    "Latency of external service calls in seconds",
			Buckets: []float64{0.005, 0.01, 0.02, 0.04},
		},
		[]string{"service", "path"},
	)

	externalServiceStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "external_service_status_total",
			Help: "Total number of responses from external services by status code",
		},
		[]string{"service", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpRequestLatency, externalServiceLatency, externalServiceStatus)
}

func MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		latency := time.Since(start)

		httpRequests.WithLabelValues(
			c.Request().Method,
			c.Path(),
			strconv.Itoa(c.Response().Status),
		).Inc()

		httpRequestLatency.WithLabelValues(
			c.Request().Method,
			c.Path(),
		).Observe(latency.Seconds())

		return err
	}
}

func RecordExternalCallMetrics(serviceName string, path string, latency time.Duration, statusCode int, callErr error) {
	externalServiceLatency.WithLabelValues(serviceName, path).Observe(latency.Seconds())
	if callErr != nil {
		externalServiceStatus.WithLabelValues(serviceName, path, "http_client_error").Inc()
	} else {
		externalServiceStatus.WithLabelValues(serviceName, path, strconv.Itoa(statusCode)).Inc()
	}
}

func AccessLogMiddleware(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			latency := time.Since(start)
			reqID := GetRequestID(c)

			logger.Infow("access log",
				"reqID", reqID,
				"method", c.Request().Method,
				"path", c.Path(),
				"status", c.Response().Status,
				"latency_ms", latency.Milliseconds(),
			)
			return err
		}
	}
}
