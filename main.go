package main

import (
	"fmt"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/plugin/storage/badger"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/ruben.vp8510/jaeger-storage-perf/generator"
	"github.com/ruben.vp8510/jaeger-storage-perf/profiling"
	"github.com/ruben.vp8510/jaeger-storage-perf/queue"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"time"
)

const nWorkers = 50
const queueSize = 2000

func EstimateNumberOfSpans(nSeconds int, writer spanstore.Writer, gen *generator.SpanGenerator) int {
	nSpans := 10000
	spans := gen.Generate(nSpans)
	q := queue.NewQueue(nWorkers, queueSize)
	q.Start(func(span model.Span) error {
		return writer.WriteSpan(&span)
	})
	start := time.Now()
	for _, span := range spans {
		q.Enqueue(span)
	}
	// Wait for all workers to stop writing
	q.Stop()
	d := time.Since(start)
	//fmt.Printf("duration=%s\n" ,d)
	//fmt.Printf("span=%f\n", float32(d.Seconds())/float32(nSpans))
	//fmt.Printf("span/s=%f\n", float32(nSpans)/float32(d.Seconds()))
	return nSeconds * int(float32(nSpans)/float32(d.Seconds()))

}

func main() {
	factory := badger.NewFactory()
	err := factory.Initialize(metrics.NullFactory, zap.NewNop())
	if err != nil {
		panic(err)
	}
	spanWriter, _ := factory.CreateSpanWriter()
	println("Estimating number of spans..")

	gen := generator.NewSpanGenerator()

	nSpans := EstimateNumberOfSpans(1, spanWriter,gen)
	spans := gen.Generate(nSpans)

	qu := queue.NewQueue(nWorkers, queueSize)

	qu.Start(func(span model.Span) error {
		return spanWriter.WriteSpan(&span)
	})

	profiler := profiling.Profiler{}

	profiler.StartProfiling()
	start := time.Now()
	for _, span := range spans {
		qu.Enqueue(span)
	}
	// Wait for all workers to stop writing
	qu.Stop()
	duration := time.Since(start)
	fmt.Printf("duration=%s\n" ,duration)
	fmt.Printf("span/s=%f\n", float32(nSpans)/float32(duration.Seconds()))
	fmt.Printf("dropped=%d\n", qu.Dropped.Load())
	profiler.StopProfiling()

}