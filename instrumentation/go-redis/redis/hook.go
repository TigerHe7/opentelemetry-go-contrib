package redis

import (
	"context"
	"go.opentelemetry.io/otel/api/trace"

	"github.com/go-redis/redis/v8"
)

type redisHook struct{
	tracer trace.Tracer
}

func (rh *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {

	newContext, _ := rh.tracer.Start(ctx, cmd.FullName()) // TODO add options support

	return newContext, nil
}

func (rh *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {

	trace.SpanFromContext(ctx).End() // TODO add options support

	return nil
}

func (rh *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	commandNamesToCounts := make(map[string]int)
	for _, cmd := range cmds {
		name := cmd.FullName()
		if _, hasKey := commandNamesToCounts[name]; hasKey {
			commandNamesToCounts[name]++
		} else {
			commandNamesToCounts[name] = 1
		}
	}
	newContext, span := rh.tracer.Start(ctx, "Begin Pipeline Process") // TODO add options support
	for k, v := range commandNamesToCounts {
		span.SetAttribute("\"" + k + "\" command count", v)
	}

	return newContext, nil
}

func (rh *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	trace.SpanFromContext(ctx).End() // TODO add options support

	return nil
}
