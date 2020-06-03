package generator

import (
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/jaegertracing/jaeger/model"
	"github.com/sethvargo/go-diceware/diceware"
	"math/rand"
)

var tagTypes = []model.ValueType{
	model.ValueType_INT64,
	model.ValueType_STRING,
	model.ValueType_BINARY,
	model.ValueType_BOOL,
	model.ValueType_FLOAT64,
}

func generateTags(minCardinality, maxCardinality int) []model.KeyValue {
	tagCard := rand.Intn(maxCardinality-minCardinality) + minCardinality
	tags := make([]model.KeyValue, tagCard)

	keys, _ := diceware.Generate(tagCard)
	for i, key := range keys {
		tags[i].Key = key
		valueType := tagTypes[rand.Intn(len(tagTypes))]
		switch valueType {
		case model.ValueType_INT64:
			tags[i].VInt64 = rand.Int63n(10000)
		case model.ValueType_FLOAT64:
			tags[i].VFloat64 = rand.Float64() * 10000
		case model.ValueType_BOOL:
			if rand.Intn(2) == 0 {
				tags[i].VBool = false
			} else {
				tags[i].VBool = true
			}
		case model.ValueType_STRING:
			v, _ := diceware.Generate(1)
			tags[i].VStr = v[0]
		}
	}
	return tags
}

func generateTraceID() model.TraceID {
	id := uuid.New()
	high := binary.LittleEndian.Uint64(id[0:8])
	low := binary.LittleEndian.Uint64(id[8:16])
	return model.TraceID{Low: low, High: high}
}

func generateRandomProcesses(num int) []*model.Process {
	processes := make([]*model.Process, num)
	names, _ := diceware.Generate(num)
	for i, srvName := range names {
		processes[i] = &model.Process{
			ServiceName: srvName,
		}
	}
	return processes
}
