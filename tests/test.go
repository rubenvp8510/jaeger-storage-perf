package tests

import (
	"context"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/rubenvp8510/jaeger-storage-perf/fixtures"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"runtime"
	"runtime/debug"
	"testing"
)

type Test struct {
	name     string
	spans    []*model.Span
	traceIds []model.TraceID
	reader   spanstore.Reader
	writer   spanstore.Writer
	factory  storage.Factory
}

func NewTest(name string) Test {
	return Test{
		name: name,
	}
}

func (t *Test) Init(fixturesPath string) error {
	println("Loading spans")
	loader := fixtures.NewLoader()
	spans, err := loader.LoadSpans(fixturesPath)
	if err != nil {
		panic(err)
	}
	t.spans = make([]*model.Span, 0, len(spans))
	traceMap := map[model.TraceID]interface{}{}
	t.traceIds = make([]model.TraceID, 0)
	for i := 0; i < len(spans); i++ {
		s := spans[i]
		if _, ok := traceMap[s.TraceID]; !ok {
			traceMap[s.TraceID] = struct{}{}
			t.traceIds = append(t.traceIds, s.TraceID)
		}
		t.spans = append(t.spans, &s)
	}
	println("Loaded.")
	runtime.GC()
	debug.FreeOSMemory()
	return nil
}

func (t *Test) InitRead() error {
	err := t.PreloadStorage(100000)
	if err != nil {
		return err
	}
	runtime.GC()
	debug.FreeOSMemory()
	return nil
}

func (t *Test) InitWrite() error {
	return nil
}

func (t *Test) ReadWriteInterface() (spanstore.Reader, spanstore.Writer) {
	return t.reader, t.writer
}

func (t *Test) SetFactory(f storage.Factory) error {
	t.factory = f
	err := f.Initialize(metrics.NullFactory, zap.NewNop())
	if err != nil {
		return err
	}
	t.writer, err = f.CreateSpanWriter()
	if err != nil {
		return err
	}

	t.reader, err = f.CreateSpanReader()
	if err != nil {
		return err
	}
	return nil
}

func (t *Test) Factory() storage.Factory {
	return t.factory
}

func (t *Test) FixtureSpans() []*model.Span {
	return t.spans
}

func (t *Test) Name() string {
	return t.name
}

func (t *Test) PreloadStorage(n int) error {
	for i := 0; i < n; i++ {
		if err := t.writer.WriteSpan(context.Background(), t.spans[i]); err != nil {
			return err
		}
	}
	return nil
}

func (t *Test) WriteBenchmark(writer spanstore.Writer, b *testing.B) (int, int) {
	dropped := 0
	count := 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := writer.WriteSpan(context.Background(), t.spans[i])
		if err != nil {
			dropped++
		} else {
			count++
		}
	}
	b.StopTimer()
	return count, dropped
}

func (t *Test) ReadBenchmark(preloadN int, reader spanstore.Reader, b *testing.B) int {
	count := 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if i >= preloadN {
			return i
		}
		_, err := reader.GetTrace(context.Background(), t.traceIds[i])
		if err != nil {
			println(err.Error())
			continue
		}
		count ++
	}
	b.StopTimer()
	return count
}
