package helpers

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/rubenvp8510/jaeger-storage-perf/generator"
	"io/ioutil"
	"math/rand"

	"context"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"io"
	"testing"
	"time"
)

var spans []*model.Span
var queries []spanstore.TraceQueryParameters

const GeneratedTraces = 1500
const MinSpans = 25
const MaxSpans = 150
const NumQueries = 1000

type SetupFunction func()storage.Factory

func saveQueryFixtures(filename string, queries []spanstore.TraceQueryParameters) error  {
	bytes, err := json.Marshal(queries)
	if err != nil {
		panic(err)
	}
	return ioutil.WriteFile(filename,bytes, 0644)
}
func saveFixtures(filename string, spans []*model.Span) error {
	batch := &model.Batch{
		Spans: spans,
	}
	bytes, err := proto.Marshal(batch)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0644)
}

func loadSpanFixtures(filename string) ([]*model.Span, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	batch := &model.Batch{}
	err = proto.Unmarshal(dat, batch)
	return batch.Spans, err
}

func loadQueriesFixtures(filename string) ([]spanstore.TraceQueryParameters, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	queries := []spanstore.TraceQueryParameters{}
	err = json.Unmarshal(dat, &queries)
	return queries, err
}
func readBenchmark(reader spanstore.Reader, b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := reader.FindTraces(context.Background(), &queries[i])
		if err != nil {
			println(err.Error())
			continue
		}
	}
	b.StopTimer()
}

func writeBenchmark(writer spanstore.Writer, b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := writer.WriteSpan(spans[i])
		if err != nil {
			println(err.Error())
		}
	}
	b.StopTimer()
}

func GenerateFixtures(filename string) error {
	println("Init generator")
	gen := generator.NewSpanGenerator()
	gen.Init()
	println("Generating traces..")
	for i := 0; i < GeneratedTraces; i++ {
		trace := gen.Generate(MinSpans, MaxSpans)
		spans = append(spans, trace...)
	}
	println("Saving fixtures..")

	queries = make([]spanstore.TraceQueryParameters, NumQueries)
	for q := 0; q < NumQueries; q++ {
		span := spans[rand.Intn(len(spans)-1)]
		opName := span.OperationName
		srvName := span.Process.ServiceName
		queries[q].OperationName = opName
		queries[q].ServiceName = srvName
		tagMap := make(map[string]string, 0)
		tagFilterNum := 0
		for i := 0; i < tagFilterNum; i++ {
			tagKey := span.Tags[rand.Intn(len(span.Tags)-1)]
			tagMap[tagKey.Key] = tagKey.AsString()
		}
		queries[q].StartTimeMin = time.Now().Add(-time.Duration(1 * time.Hour))
		queries[q].StartTimeMax = time.Now().Add(time.Duration(1 * time.Hour))
	}
	err := saveQueryFixtures(filename+"_queries", queries)
	if err != nil {
		return err
	}
	return saveFixtures(filename, spans)
}

func Setup(fixturesPath string) {
	println("Loading fixtures..")
	var err error
	spans, err = loadSpanFixtures(fixturesPath)
	if err != nil {
		panic(err)
	}
	queries , err = loadQueriesFixtures(fixturesPath+"_queries")
	if err != nil {
		panic(err)
	}
	println("Fixtures loaded..")


}

func StorageBenchmark(stype string, factory storage.Factory, b *testing.B) {
	b.StopTimer()
	err := factory.Initialize(metrics.NullFactory, zap.NewNop())
	if err != nil {
		panic(err)
	}
	spanWriter, _ := factory.CreateSpanWriter()
	spanReader, _ := factory.CreateSpanReader()

	b.Run(stype+"-write", func(b *testing.B) {
		writeBenchmark(spanWriter, b)
	})

	b.Run(stype+"-read", func(b *testing.B) {
		readBenchmark(spanReader, b)
	})

	if closer, ok := factory.(io.Closer); ok {
		closer.Close()
	}
}
