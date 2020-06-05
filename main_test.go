package main

import (
	"github.com/jaegertracing/jaeger/plugin/storage/badger"
	"github.com/ruben.vp8510/jaeger-storage-perf/generator"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"testing"
)



func BenchmarkBadgerStorage(b *testing.B) {
	b.StopTimer()
	factory := badger.NewFactory()
	err := factory.Initialize(metrics.NullFactory, zap.NewNop())
	if err != nil {
		panic(err)
	}
	spanWriter, _ := factory.CreateSpanWriter()
	gen := generator.NewSpanGenerator()
	spanList := gen.Generate(b.N+1)
	b.StartTimer()
	for i:= 0; i < b.N ; i++ {
		spanWriter.WriteSpan(&spanList[i])
	}
	b.StopTimer()
}
