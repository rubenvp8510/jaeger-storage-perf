package perftest

import (
	"github.com/jaegertracing/jaeger/storage"
	"github.com/rubenvp8510/jaeger-storage-perf/tests"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func runWrite(test StorageTest, b *testing.B) error {
	spanReader, spanWriter := test.ReadWriteInterface()
	b.StopTimer()
	sentSpans, droppedSpans := test.WriteBenchmark(spanWriter, b)
	writtenSpans := verifyWrites(test.FixtureSpans(), sentSpans, spanReader)
	b.ReportMetric(float64(sentSpans), "sent")
	b.ReportMetric(float64(writtenSpans), "verified")
	b.ReportMetric(float64(droppedSpans), "dropped")
	return closeFactory(test)
}

func runRead(preloadN int, test StorageTest, b *testing.B) error {
	spanReader, _ := test.ReadWriteInterface()
	readSpans := test.ReadBenchmark(preloadN, spanReader, b)
	b.ReportMetric(float64(readSpans), "read")
	return closeFactory(test)
}

func RunRead(t *testing.B, fixturesPath string, factories map[string]storage.Factory) {
	test := tests.NewTest(t.Name())
	err := test.Init(fixturesPath)
	assert.NoError(t, err)
	for storageType, factory := range factories {
		err := test.SetFactory(factory)
		assert.NoError(t, err)
		err = test.InitRead()
		assert.NoError(t, err)
		t.Run(storageType, func(b *testing.B) {
			b.StopTimer()
			err := runRead(100000, &test, b)
			assert.NoError(t, err)
		})
	}
}

func RunWrite(t *testing.B, fixturesPath string, factories map[string]storage.Factory) {
	test := tests.NewTest(t.Name())
	err := test.Init(fixturesPath)
	assert.NoError(t, err)
	for storageType, factory := range factories {
		err := test.SetFactory(factory)
		assert.NoError(t, err)
		err = test.InitWrite()
		assert.NoError(t, err)
		t.Run(storageType, func(b *testing.B) {
			b.StopTimer()
			err := runWrite(&test, b)
			assert.NoError(t, err)
		})
	}
}
func closeFactory(test StorageTest) error {
	if closer, ok := test.Factory().(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
