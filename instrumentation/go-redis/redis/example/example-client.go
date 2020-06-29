// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	redisWrapper "go.opentelemetry.io/contrib/instrumentation/go-redis/redis"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"os"
	"time"
)

func main() {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
	})
	defer rdb.Close()

	logger := log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile)
	exporter, err := zipkin.NewExporter(
		"http://localhost:9411/api/v2/spans",
		"zipkin-example",
		zipkin.WithLogger(logger),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5),
			sdktrace.WithMaxExportBatchSize(10),
		),
	)

	ctx := context.Background()
	tracer := tp.Tracer("Example tracer")
	wrappedRdb := redisWrapper.Wrap(rdb).WithContext(ctx, tracer)

	pipeline := wrappedRdb.Pipeline()
	_, err = pipeline.Pipelined(ctx, func (pipeliner redis.Pipeliner) error {
		pipeliner.Set(ctx, "rand key", "val", 0)
		//<- time.After(2 * time.Second)
		log.Print(pipeliner.Get(ctx, "rand key").Result())
		return nil
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	//for i := 1; i < 10; i++ {
	//	wrappedRdb.Set(ctx, fmt.Sprintf("example_int_key_%d", i), i, 1000000000)
	//	wrappedRdb.Set(ctx, fmt.Sprintf("example_string_key_%d", i), "example_value", 1000000000)
	//}
	//
	//intVal := wrappedRdb.Get(ctx, "example_int_key_5").String()
	//stringVal := wrappedRdb.Get(ctx, "example_string_key_5").String()
	<- time.After(2 * time.Second)
	//fmt.Printf("%s", intVal)
	//fmt.Printf("%s", stringVal)
}
