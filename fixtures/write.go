package fixtures

import (
	"encoding/gob"
	"encoding/json"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/rubenvp8510/jaeger-storage-perf/generator"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"time"
)

func SaveSpans(filename string, nTraces, minSpans, maxSpans int) (int, error) {
	os.Remove(filename)
	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	enc := gob.NewEncoder(file)
	defer func() {
		_ = file.Close()
	}()
	spans := make([]*model.Span,0,nTraces*minSpans)
	gen := generator.NewSpanGenerator()
	gen.Init()
	for i := 0; i < nTraces; i++ {
		spans = append(spans, gen.Generate(minSpans, maxSpans)...)
	}
	err := enc.Encode(spans)
	if err != nil {
		return 0, err
	}
	return len(spans) , nil
}

func SaveQueries(filename string, queries []spanstore.TraceQueryParameters) error {
	bytes, err := json.Marshal(queries)
	if err != nil {
		panic(err)
	}
	return ioutil.WriteFile(filename, bytes, 0644)
}


// TODO: Move this to the generator
func GenerateQuery(spans []*model.Span, numQueries int) []spanstore.TraceQueryParameters {
	type kv struct {
		Key     string
		Value   string
		Matches int
	}

	tags := make(map[string]kv)
	for _, span := range spans {
		for _, tag := range span.Tags {
			key := tag.Key + tag.AsString()
			if _, ok := tags[key]; ok {
				matches := tags[key]
				matches.Matches++
				tags[key] = matches
			} else {
				tags[key] = kv{
					Key:     tag.Key,
					Value:   tag.AsString(),
					Matches: 1,
				}
			}
		}
	}
	sortedByMatch := make([]kv, len(tags))
	count := 0
	for _, item := range tags {
		sortedByMatch[count] = item
		count++
	}
	sort.Slice(sortedByMatch, func(i, j int) bool {
		return sortedByMatch[i].Matches > sortedByMatch[j].Matches
	})

	queriableTags := make([]kv, 0)

	for _, t := range sortedByMatch {
		if t.Matches > 1 {
			queriableTags = append(queriableTags, t)
		}
	}
	queries := make([]spanstore.TraceQueryParameters, numQueries)

	for q := 0; q < numQueries; q++ {
		span := spans[rand.Intn(len(spans)-1)]
		queries[q].ServiceName = span.Process.ServiceName

		// Filter by operation
		if rand.Intn(2) == 0 {
			queries[q].OperationName = span.OperationName
		}

		tagMap := make(map[string]string, 0)
		tagFilterNum := rand.Intn(3)
		for i := 0; i < tagFilterNum; i++ {
			tag := queriableTags[rand.Intn(len(queriableTags)-1)]
			tagMap[tag.Key] = tag.Value
		}
		queries[q].Tags = tagMap
		queries[q].StartTimeMin = time.Now().Add(-time.Duration(1 * time.Hour))
		queries[q].StartTimeMax = time.Now().Add(time.Duration(1 * time.Hour))
	}
	return queries

}
