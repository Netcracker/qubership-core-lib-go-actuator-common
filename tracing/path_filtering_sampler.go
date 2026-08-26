package tracing

import (
	"strings"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

// DefaultExcludedTracePaths are exact HTTP paths dropped from sampling by default
// (query string is stripped before matching). Override via ZipkinOptions.ExcludedPaths.
var DefaultExcludedTracePaths = []string{
	"/health",
	"/ready",
	"/livez",
	"/readyz",
	"/healthz",
	"/liveness",
	"/readiness",
	"/prometheus",
	"/metrics",
	"/api-version",
}

// DefaultExcludedTracePathPrefixes are HTTP path prefixes dropped from sampling by default.
// Override via ZipkinOptions.ExcludedPathPrefixes.
var DefaultExcludedTracePathPrefixes = []string{
	"/static",
}

// PathFilteringSampler drops spans whose name or http.target matches excluded paths/prefixes,
// then delegates the decision to the wrapped sampler.
type PathFilteringSampler struct {
	delegate sdktrace.Sampler
	paths    map[string]struct{}
	prefixes []string
}

// NewPathFilteringSampler wraps delegate with path-based drop rules.
// paths are exact matches; prefixes use strings.HasPrefix after query stripping.
// Empty paths and prefixes means no filtering (delegate always used).
func NewPathFilteringSampler(delegate sdktrace.Sampler, paths []string, prefixes []string) sdktrace.Sampler {
	pathSet := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		if p == "" {
			continue
		}
		pathSet[p] = struct{}{}
	}
	copiedPrefixes := make([]string, 0, len(prefixes))
	for _, p := range prefixes {
		if p == "" {
			continue
		}
		copiedPrefixes = append(copiedPrefixes, p)
	}
	return PathFilteringSampler{
		delegate: delegate,
		paths:    pathSet,
		prefixes: copiedPrefixes,
	}
}

func (s PathFilteringSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	if s.isExcluded(p.Name) || s.isExcluded(httpTargetFromSamplingParams(p)) {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.Drop,
			Tracestate: psc.TraceState(),
		}
	}
	return s.delegate.ShouldSample(p)
}

func (s PathFilteringSampler) Description() string {
	return "PathFilteringSampler{" + s.delegate.Description() + "}"
}

func (s PathFilteringSampler) isExcluded(path string) bool {
	if path == "" {
		return false
	}
	path = stripQuery(path)
	if _, ok := s.paths[path]; ok {
		return true
	}
	for _, prefix := range s.prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func stripQuery(path string) string {
	if i := strings.IndexByte(path, '?'); i >= 0 {
		return path[:i]
	}
	return path
}

func httpTargetFromSamplingParams(p sdktrace.SamplingParameters) string {
	for _, attr := range p.Attributes {
		if attr.Key == semconv.HTTPTargetKey {
			if val := attr.Value.AsString(); val != "" {
				return val
			}
		}
	}
	return ""
}

// ResolveExcludedTracePaths returns the effective exclusion lists for ZipkinOptions.
// A nil slice means "use library defaults"; a non-nil (including empty) slice replaces defaults.
func ResolveExcludedTracePaths(opts *ZipkinOptions) (paths []string, prefixes []string) {
	if opts == nil {
		return append([]string(nil), DefaultExcludedTracePaths...), append([]string(nil), DefaultExcludedTracePathPrefixes...)
	}
	if opts.ExcludedPaths == nil {
		paths = append([]string(nil), DefaultExcludedTracePaths...)
	} else {
		paths = append([]string(nil), opts.ExcludedPaths...)
	}
	if opts.ExcludedPathPrefixes == nil {
		prefixes = append([]string(nil), DefaultExcludedTracePathPrefixes...)
	} else {
		prefixes = append([]string(nil), opts.ExcludedPathPrefixes...)
	}
	return paths, prefixes
}
