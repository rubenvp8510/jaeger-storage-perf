package generator

import (
	"github.com/google/uuid"
	"github.com/jaegertracing/jaeger/model"
	"github.com/sethvargo/go-diceware/diceware"
	"math/rand"
	"time"
)

type SpanGenerator struct {
	minNumberTagsSpan    int
	maxNumberTagsSpan    int
	minNumberTagsProcess int
	maxNumberTagsProcess int
}

func NewSpanGenerator() *SpanGenerator {
	return &SpanGenerator{}
}

func (g *SpanGenerator) generateSpan(traceID model.TraceID, traceStartTime, traceEndTime int64) model.Span {
	processes := generateRandomProcesses(10)
	opNames, _ := diceware.Generate(1)
	startTime := rand.Int63n(traceEndTime-traceStartTime) + traceStartTime
	duration := rand.Int63n(traceEndTime - startTime)

	return model.Span{
		TraceID:       traceID,
		SpanID:        model.SpanID(rand.Int63()),
		OperationName: opNames[0],
		Tags:          generateTags(0, 100),
		Process:       processes[rand.Intn(len(processes))],
		StartTime:     time.Unix(startTime, 0),
		Duration:      time.Duration(duration) * time.Second,
		Flags:         model.Flags(rand.Int31n(8)),
		ProcessID:     uuid.New().String(),
	}
}

func (g *SpanGenerator) Generate(numSpans int) []model.Span {
	traceID := generateTraceID()
	duration := rand.Int63n(9) + 1
	timestamp := time.Now().Unix() - rand.Int63n(1000)
	var spans []model.Span
	for i := 0; i < numSpans; i++ {
		span := g.generateSpan(traceID, timestamp, timestamp+duration)
		spans = append(spans, span)
	}
	return spans
}
