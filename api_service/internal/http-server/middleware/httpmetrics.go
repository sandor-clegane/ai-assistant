package middleware

import (
	"ai-assistant-api/internal/metrics"
	helper "ai-assistant-api/internal/utils/response-writer"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func HTTPMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.HttpRequestsCurrent.WithLabelValues().Inc()

		srw := helper.NewStatusResponseWriter(w)
		now := time.Now()

		next.ServeHTTP(srw, r)

		elapsedSeocnds := time.Since(now).Seconds()
		pattern := chi.RouteContext(r.Context()).RoutePattern()
		method := chi.RouteContext(r.Context()).RouteMethod
		status := srw.GetStatusString()

		metrics.HttpRequestsCurrent.WithLabelValues().Dec()
		metrics.HttpRequestsTotal.WithLabelValues(pattern, method, status).Inc()
		metrics.HttpRequestsDurationHistorgram.WithLabelValues(pattern, method).Observe(elapsedSeocnds)
	})
}
