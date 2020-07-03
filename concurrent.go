package main

import (
	"github.com/jaegertracing/jaeger/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rubenvp8510/jaeger-storage-perf/generator"
	"github.com/rubenvp8510/jaeger-storage-perf/profiling"
	"github.com/rubenvp8510/jaeger-storage-perf/queue"
	"github.com/rubenvp8510/jaeger-storages/questbd"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const nWorkers = 1
const queueSize = 10

func performanceTest(stop chan struct{}) {
	go func() {
		println("Starting test...")
		//factory :=
		// .NewFactory()
		// factory := redbull.NewFactory()
		factory := questbd.NewFactory()
		/*factory.Options.GetPrimary().Ephemeral = false
		factory.Options.GetPrimary().KeyDirectory = "/home/rvargasp/badger/"
		factory.Options.GetPrimary().ValueDirectory = "/home/rvargasp/badger/"
		factory.Options.GetPrimary().SyncWrites = true*/

		err := factory.Initialize(metrics.NullFactory, zap.NewNop())
		// defer factory.Close()

		if err != nil {
			panic(err)
		}
		spanWriter, _ := factory.CreateSpanWriter()
		qu := queue.NewQueue(nWorkers, queueSize)
		// Starting prometheus reporter
		reporter := profiling.NewReporter(qu)
		reporter.Start()
		// Starting trace generator.
		gen := generator.NewSpanGenerator()
		gen.Init()

		// Start queue producer-consumer
		qu.Start(func(span model.Span) error {
			err := spanWriter.WriteSpan(&span)
			if err != nil {
				return err
			}
			return nil
		})

		// Buffer of traces
		traces := gen.Generate(50,500000)
		println(len(traces))
		for _, trace := range traces {
			qu.Enqueue(trace)
		}
		<- stop
		println("Finish")

	}()
}

func main() {
	sigs := make(chan os.Signal, 1)
	stop := make(chan struct{}, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	performanceTest(stop)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

	<-sigs
	stop <- struct{}{}


}
