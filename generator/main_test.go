package generator

import (
	"testing"
)

func BenchmarkTraceGenerator_Generate(b *testing.B) {
	gen := NewSpanGenerator()
	gen.Init()
	b.StopTimer()
	b.StartTimer()
	for i:= 0; i < b.N ; i++ {
		gen.Generate(50,51)
	}
	b.StopTimer()
}


