package fixtures

import "testing"

func TestLoading(t *testing.T) {
	loader := NewLoader()
	traces, _ := loader.LoadSpans("/home/rvargasp/go/src/github.com/rubenvp8510/jaeger-storage-perf/traces")
	println(len(traces))
}

func BenchmarkLoadFixtures(b *testing.B) {
	b.ReportAllocs()

	loader := NewLoader()
	b.ResetTimer()
	for i:=0; i < b.N; i++ {
		loader.LoadSpans("/home/rvargasp/go/src/github.com/rubenvp8510/jaeger-storage-perf/traces")
	}
	b.StopTimer()
}