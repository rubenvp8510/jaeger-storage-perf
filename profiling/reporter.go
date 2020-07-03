package profiling

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rubenvp8510/jaeger-storage-perf/queue"
	"time"
)

type Reporter struct {
	q                   *queue.Queue
	cumulativeSpans     int64
	spansPerSecond      int64
	spansPerSecondGauge prometheus.Gauge
	totalSpansCounter   prometheus.Counter
}

func NewReporter(q *queue.Queue) *Reporter {
	totalSpans := promauto.NewCounter(prometheus.CounterOpts{
		Name: "jaeger_perf_spans_total",
		Help: "The total number of stored spans",
	})
	spansPerSecond := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "jaeger_perf_spans_second",
		Help: "Spans per second",
	})
	return &Reporter{
		q:                   q,
		spansPerSecondGauge: spansPerSecond,
		totalSpansCounter:   totalSpans,
		cumulativeSpans:     0,
		spansPerSecond:      0,
	}
}

func (s *Reporter) Start() {
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				newAcum := int64(s.q.Added.Load())
				oldAcum := s.cumulativeSpans
				speed := newAcum - oldAcum
				deltaGauge := speed - s.spansPerSecond
				s.spansPerSecond = speed
				s.cumulativeSpans = newAcum
				s.totalSpansCounter.Add(float64(speed))
				s.spansPerSecondGauge.Add(float64(deltaGauge))
				dropped := s.q.Dropped.Load()
				fmt.Printf("%d spans/sec dropped: %d \n", s.spansPerSecond, dropped)
			}
		}
	}()
}
