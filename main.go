package main

import (
	"fmt"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/plugin/storage/badger"
	"github.com/ruben.vp8510/jaeger-storage-perf/generator"
	"github.com/ruben.vp8510/jaeger-storage-perf/profiling"
	"github.com/ruben.vp8510/jaeger-storage-perf/queue"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"time"
)

func main() {
	factory := badger.NewFactory()
	err := factory.Initialize(metrics.NullFactory, zap.NewNop())
	if err != nil {
		panic(err)
	}
	spanWriter, _ := factory.CreateSpanWriter()
	println("Generating test data..")
	gen := generator.NewSpanGenerator()
	spans := gen.Generate(10000)

	profiler := profiling.Profiler{}
	q := queue.NewQueue(10, 100)
	q.Start(func(span model.Span) error {
		return spanWriter.WriteSpan(&span)
	})

	profiler.StartProfiling()
	println("Starting performance test..")
	start := time.Now()
	for _, span := range spans {
		q.Enqueue(span)
	}
	// Wait for all workers to stop writing
	q.Stop()
	fmt.Printf("duration=%s\n" ,time.Since(start))
	profiler.StopProfiling()

}
