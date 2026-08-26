package tracing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func TestPathFilteringSamplerDropsDefaultProbesAndManagement(t *testing.T) {
	paths, prefixes := ResolveExcludedTracePaths(&ZipkinOptions{})
	sampler := NewPathFilteringSampler(NewRateLimitingSampler(10), paths, prefixes)

	for _, path := range []string{"/health", "/ready", "/prometheus", "/api-version", "/metrics", "/ready?x=1", "/static/app.js"} {
		decision := sampler.ShouldSample(sdktrace.SamplingParameters{
			Name: "/",
			Attributes: []attribute.KeyValue{
				semconv.HTTPTargetKey.String(path),
			},
		}).Decision
		assert.Equal(t, sdktrace.Drop, decision, "expected drop for %s", path)
	}

	decision := sampler.ShouldSample(sdktrace.SamplingParameters{
		Name: "/ready",
	}).Decision
	assert.Equal(t, sdktrace.Drop, decision)
}

func TestPathFilteringSamplerKeepsBusinessPaths(t *testing.T) {
	paths, prefixes := ResolveExcludedTracePaths(&ZipkinOptions{})
	sampler := NewPathFilteringSampler(NewRateLimitingSampler(10), paths, prefixes)

	decision := sampler.ShouldSample(sdktrace.SamplingParameters{
		Name: "/api/v1/something",
		Attributes: []attribute.KeyValue{
			semconv.HTTPTargetKey.String("/api/v1/something"),
		},
	}).Decision
	assert.Equal(t, sdktrace.RecordAndSample, decision)
}

func TestPathFilteringSamplerCustomPathsOverrideDefaults(t *testing.T) {
	sampler := NewPathFilteringSampler(
		NewRateLimitingSampler(10),
		[]string{"/custom-probe"},
		[]string{},
	)

	assert.Equal(t, sdktrace.Drop, sampler.ShouldSample(sdktrace.SamplingParameters{
		Attributes: []attribute.KeyValue{semconv.HTTPTargetKey.String("/custom-probe")},
	}).Decision)

	// defaults are not applied when overridden
	assert.Equal(t, sdktrace.RecordAndSample, sampler.ShouldSample(sdktrace.SamplingParameters{
		Attributes: []attribute.KeyValue{semconv.HTTPTargetKey.String("/health")},
	}).Decision)
}

func TestPathFilteringSamplerEmptyListsDisableFiltering(t *testing.T) {
	sampler := NewPathFilteringSampler(NewRateLimitingSampler(10), []string{}, []string{})
	assert.Equal(t, sdktrace.RecordAndSample, sampler.ShouldSample(sdktrace.SamplingParameters{
		Attributes: []attribute.KeyValue{semconv.HTTPTargetKey.String("/health")},
	}).Decision)
}

func TestPathFilteringSamplerDescription(t *testing.T) {
	sampler := NewPathFilteringSampler(NewRateLimitingSampler(10), DefaultExcludedTracePaths, nil)
	assert.Contains(t, sampler.Description(), "PathFilteringSampler")
	assert.Contains(t, sampler.Description(), "RateLimitingSampler")
}

func TestResolveExcludedTracePathsNilMeansDefaults(t *testing.T) {
	paths, prefixes := ResolveExcludedTracePaths(&ZipkinOptions{})
	assert.Equal(t, DefaultExcludedTracePaths, paths)
	assert.Equal(t, DefaultExcludedTracePathPrefixes, prefixes)
}

func TestResolveExcludedTracePathsEmptyOverridesDefaults(t *testing.T) {
	paths, prefixes := ResolveExcludedTracePaths(&ZipkinOptions{
		ExcludedPaths:        []string{},
		ExcludedPathPrefixes: []string{},
	})
	assert.Empty(t, paths)
	assert.Empty(t, prefixes)
}

func TestHttpTargetFromSamplingParamsIgnoresEmptyAndMissing(t *testing.T) {
	assert.Equal(t, "", httpTargetFromSamplingParams(sdktrace.SamplingParameters{}))
	assert.Equal(t, "", httpTargetFromSamplingParams(sdktrace.SamplingParameters{
		Attributes: []attribute.KeyValue{
			semconv.HTTPTargetKey.String(""),
			semconv.HTTPMethodKey.String("GET"),
		},
	}))
	assert.Equal(t, "/ready", httpTargetFromSamplingParams(sdktrace.SamplingParameters{
		Attributes: []attribute.KeyValue{
			semconv.HTTPMethodKey.String("GET"),
			semconv.HTTPTargetKey.String("/ready"),
		},
	}))
}
