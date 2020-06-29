module go.opentelemetry.io/contrib/instrumentation/go-redis/redis

go 1.13

require (
	github.com/go-redis/redis/v8 v8.0.0-beta.5
	go.opentelemetry.io/otel v0.6.0
	go.opentelemetry.io/otel/exporters/trace/zipkin v0.6.0
)
