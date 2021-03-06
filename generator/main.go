package generator

import (
	"github.com/jaegertracing/jaeger/model"
	"math/rand"
	"time"
)

type TraceGenerator struct {
	MaxTags        int
	MinTags        int
	MaxProcess     int
	MinProcessTags int
	MaxDuration    time.Duration
	MinDuration    time.Duration

	processes []Process
	tags      []*TagTemplate
	ops       []string
}

func NewSpanGenerator() *TraceGenerator {
	return &TraceGenerator{
		MaxProcess:  10,
		MaxTags:     100,
		MinTags:     2,
		MaxDuration: time.Duration(10 * time.Second),
		MinDuration: time.Duration(1 * time.Second),
	}
}

func (g *TraceGenerator) Init() {
	rand.Seed(time.Now().Unix())
	words := generateWords(20000)
	g.tags = generateTagTemplates(g.MaxTags, words)
	g.processes = generateProcesses(g.MaxProcess+1, g.MinProcessTags, g.tags)
	ops := generateWords(g.MaxProcess + 1)
	g.ops = ops

}

func (g *TraceGenerator) generateTags(min int) []model.KeyValue {
	return generateTagsFromPool(g.tags, min)
}

func (g *TraceGenerator) generateSpan(traceID model.TraceID, traceStartTime, traceEndTime int64) *model.Span {
	opNames := g.ops[rand.Intn(g.MaxProcess-1)]
	startTime := rand.Int63n(traceEndTime-traceStartTime) + traceStartTime
	duration := rand.Int63n(traceEndTime-startTime) + 1
	process := g.processes[rand.Intn(g.MaxProcess)]
	return &model.Span{
		TraceID:       traceID,
		SpanID:        model.SpanID(rand.Int63()),
		OperationName: opNames,
		Tags:          g.generateTags(g.MinTags),
		Process:       process.Process,
		StartTime:     time.Unix(startTime, 0),
		Duration:      time.Duration(duration) * time.Second,
		Flags:         model.Flags(0),
		ProcessID:     process.Id,
	}
}

func (g *TraceGenerator) setRelations(traceID model.TraceID, spans []*model.Span, desiredLevels int) {
	rand.Shuffle(len(spans), func(i, j int) { spans[i], spans[j] = spans[j], spans[i] })
	avgSpanNum := (len(spans) - 1) / desiredLevels
	levels := 0
	lowerParent := 0
	upperParent := 1
	upperChild := upperParent + avgSpanNum

	for {
		if levels >= desiredLevels {
			break
		}
		if upperParent >= len(spans) {
			upperParent = len(spans) - 1
		}
		if upperChild >= len(spans) {
			upperChild = len(spans) - 1
		}
		if upperParent == upperChild {
			break
		}
		pool := spans[lowerParent:upperParent]
		children := spans[upperParent:upperChild]
		maxParent := 0
		if len(pool) > 1 {
			maxParent = rand.Intn(len(pool)) - 1
		}
		for i := 0; i < len(children); i++ {
			parentIndex := 0
			if maxParent > 1 {
				parentIndex = rand.Intn(maxParent)
			}
			children[i].References = append(children[i].References, model.SpanRef{
				TraceID: traceID,
				SpanID:  pool[parentIndex].SpanID,
				RefType: model.SpanRefType_CHILD_OF,
			})
		}
		lowerParent = upperParent
		upperParent = upperParent + avgSpanNum
		upperChild = upperParent + avgSpanNum
		levels++
	}
}

func (g *TraceGenerator) Generate(minSpans, maxSpans int) []*model.Span {
	numSpans := generateRandomNumber(minSpans, maxSpans)
	traceID := generateTraceID()
	duration := generateRandomNumberInt64(int(g.MinDuration.Seconds()), int(g.MaxDuration.Seconds()))
	timestamp := time.Now().Unix() - rand.Int63n(1000)
	var spans []*model.Span
	for i := 0; i < numSpans; i++ {
		span := g.generateSpan(traceID, timestamp, timestamp+duration)
		spans = append(spans, span)
	}
	g.setRelations(traceID, spans, 5)
	result := make([]*model.Span, numSpans)
	for i, s := range spans {
		result[i] = s
	}
	return result
}
