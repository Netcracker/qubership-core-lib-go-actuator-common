module github.com/netcracker/qubership-core-lib-go-actuator-common/v2

go 1.27.1

retract (
	v2.10.0 // Published by error.
	v2.9.1 // Published by error.
	v2.9.0 // Published by error.
	v2.8.0 // Published by error.
	v2.7.4 // Published by error.
	v2.7.3 // Published by error.
	v2.7.2 // Published by error.
	v2.7.1 // Published by error.
)

require (
	github.com/netcracker/qubership-core-lib-go/v3 v3.14.0
	github.com/prometheus/client_golang v1.24.1
	github.com/prometheus/client_model v0.6.2
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/otel v1.46.0
	go.opentelemetry.io/otel/exporters/zipkin v1.46.0
	go.opentelemetry.io/otel/sdk v1.46.0
	go.opentelemetry.io/otel/trace v1.46.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/knadh/koanf/maps v0.1.3 // indirect
	github.com/knadh/koanf/parsers/yaml v1.1.1 // indirect
	github.com/knadh/koanf/providers/env/v2 v2.0.1 // indirect
	github.com/knadh/koanf/providers/file v1.2.1 // indirect
	github.com/knadh/koanf/v2 v2.3.6 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	github.com/viney-shih/go-lock v1.1.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/metric v1.46.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
