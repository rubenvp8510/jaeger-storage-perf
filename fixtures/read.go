package fixtures

import (
	"bufio"
	"encoding/gob"
	"github.com/jaegertracing/jaeger/model"
	"os"
)

const BufferSize = 2 * 1024


type Loader struct {}

func NewLoader() *Loader {
	loader := &Loader{}
	return loader
}

func (l *Loader) LoadSpans(filename string) ([]model.Span, error) {
	file, _ := os.OpenFile(filename, os.O_RDONLY, 0644)
	reader := bufio.NewReaderSize(file, BufferSize)
	dec := gob.NewDecoder(reader) // Will read from network.
	spans := make([]model.Span, 0, 10)
	err := dec.Decode(&spans)
	if err != nil {
		return nil, err
	}
	return spans, nil
}