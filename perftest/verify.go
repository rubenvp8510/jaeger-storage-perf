package perftest

import (
	"context"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

type SetupFunction func() storage.Factory

func verifyWrites(spans []*model.Span, count int, reader spanstore.Reader) int {
	verified := 0
	traces := make(map[string]struct{})
	spansDB := make(map[string]*model.Span)
	for i := 0; i < count; i++ {
		span := spans[i]
		traceId := span.TraceID
		if _, ok := traces[traceId.String()]; !ok {
			trace, err := reader.GetTrace(context.Background(), traceId)
			if err == nil && trace != nil {
				traces[traceId.String()] = struct{}{}
				for _, s := range trace.Spans {
					spansDB[s.SpanID.String()] = s
				}
			}
			if _, hasSpan := spansDB[span.SpanID.String()]; hasSpan {
				verified++
			}
		} else {
			if _, hasSpan := spansDB[span.SpanID.String()]; hasSpan {
				verified++
			}
		}
	}
	return verified
}

