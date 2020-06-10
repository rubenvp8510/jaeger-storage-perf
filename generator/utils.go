package generator

import (
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/jaegertracing/jaeger/model"
	"github.com/ruben.vp8510/jaeger-storage-perf/generator/data"
)

func generateTraceID() model.TraceID {
	id := uuid.New()
	traceID := model.TraceID{}
	traceID.High = binary.BigEndian.Uint64(id[:8])
	traceID.Low = binary.BigEndian.Uint64(id[8:])
	return traceID
}

func generateProcesses(num int, minTags int, template []*TagTemplate) []Process {
	processes := make([]Process, num)
	names := generateWords(num)
	for i, srvName := range names {
		processes[i].Process = &model.Process{
			ServiceName: srvName,
			Tags:generateTagsFromPool(template,minTags),
		}
		processes[i].Id = uuid.New().String()
	}
	return processes
}

func generateWords(max int) []string  {
	return generateRandStrings(data.Words, max)
}


func generateRandStrings(pool []string, max int) []string {
	size := len(pool)
	tagKeys := make([]string, max)
	count := 0

	for i := 0; i < size && count < max; i, count = i+1, count+1 {
		tagKeys[count] = pool[i]
	}

	for {
		m := count
		for i := 0; i < m; i++ {
			prefix := tagKeys[i]
			for k := 0; k < size; k++ {
				key := prefix + tagSeparator + pool[k]
				count++
				if count >= max {
					return tagKeys
				}
				tagKeys[count] = key
			}
		}
	}
}