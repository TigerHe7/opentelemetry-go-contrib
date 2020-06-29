package redis

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"

	"github.com/go-redis/redis/v8"
)

// WrappedClient is the interface returned by Wrap.
//
// WrappedClient implements redis.UniversalClient
type WrappedClient interface {
	redis.UniversalClient
}

type clientWithContext struct {
	*redis.Client
}

// Wrap wraps client such that executed commands are reported as spans to Elastic APM,
// using the client's associated context.
// A context-specific client may be obtained by using Client.WithContext.
func Wrap(client redis.UniversalClient) *clientWithContext { // TODO: pass in tracer config or own tracer
	switch client.(type) {
	case *redis.Client:
		return &clientWithContext{Client: client.(*redis.Client)}
	}

	return nil // client.(WrappedClient)
}

func (c clientWithContext) WithContext(ctx context.Context, tracer trace.Tracer) WrappedClient {
	c.Client = c.Client.WithContext(ctx)

	c.Client.AddHook(&redisHook{tracer})
	return c
}
