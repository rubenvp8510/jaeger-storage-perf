package perftest

import (
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"testing"
)

type StorageTest interface {
	Name() string
	Init(fixturesPath string) error
	InitRead() error
	InitWrite() error
	SetFactory(f storage.Factory) error
	Factory() storage.Factory
	ReadWriteInterface() (spanstore.Reader, spanstore.Writer)
	FixtureSpans() []*model.Span
	PreloadStorage(n int) error
	WriteBenchmark(writer spanstore.Writer, b *testing.B) (int, int)
	ReadBenchmark(preloadN int, reader spanstore.Reader, b *testing.B) int
}
