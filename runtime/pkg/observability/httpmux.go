package observability

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// MuxHandle is a wrapper around http.ServeMux.Handle that adds route tags to the handler.
// It does NOT wrap the handler with observability.Middleware. The caller is expected to add that on the ServeMux itself.
//
// It reimplements otelhttp.WithRouteTag, which was removed in otelhttp v0.69.0, by tagging the active span and the
// request labeler with the http.route attribute, exactly as WithRouteTag did.
func MuxHandle(mux *http.ServeMux, pattern string, handler http.Handler) {
	attr := semconv.HTTPRouteKey.String(pattern)
	mux.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trace.SpanFromContext(r.Context()).SetAttributes(attr)
		if labeler, ok := otelhttp.LabelerFromContext(r.Context()); ok {
			labeler.Add(attr)
		}
		handler.ServeHTTP(w, r)
	}))
}
