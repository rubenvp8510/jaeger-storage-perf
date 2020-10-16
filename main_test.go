package main

import (
	"github.com/rubenvp8510/jaeger-storage-perf/perftest"
	"os"
	"testing"
)

func BenchmarkBlobStorageWrite(t *testing.B) {
	perftest.RunWrite(t)
}


func BenchmarkBlobStorageRead(t *testing.B) {
	perftest.RunRead(t)
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
