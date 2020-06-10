package generator

import (
	"github.com/jaegertracing/jaeger/model"
	"github.com/ruben.vp8510/jaeger-storage-perf/generator/data"
	"math"
	"math/rand"
)

const tagSeparator = "."

var tagTypes = []model.ValueType{
	model.ValueType_INT64,
	model.ValueType_STRING,
	model.ValueType_BOOL,
	model.ValueType_FLOAT64,
}

type Process struct {
	Id      string
	Process *model.Process
}

type TagTemplate struct {
	Key  string
	Type model.ValueType
	words []string
}

func (t *TagTemplate) Tag() model.KeyValue {
	tag := model.KeyValue{
		Key:   t.Key,
		VType: t.Type,
	}
	switch t.Type {
	case model.ValueType_INT64:
		tag.VInt64 = rand.Int63n(math.MaxInt64)
	case model.ValueType_FLOAT64:
		tag.VFloat64 = rand.Float64() * math.MaxFloat64
	case model.ValueType_BOOL:
		tag.VBool = rand.Intn(2) == 0
	case model.ValueType_STRING:
		tag.VStr = t.words[rand.Intn(len(t.words) - 1)]
	}
	return tag
}

func generateTagTemplates(max int, words[]string) []*TagTemplate {
	tags := make([]*TagTemplate, max)
	keys := generateRandStrings(data.Tags,tagSeparator, max)
	ntypes := len(tagTypes) - 1
	for i := 0; i < len(tags); i++ {
		t := rand.Intn(ntypes)
		tags[i] = &TagTemplate{
			Key:  keys[i],
			Type: tagTypes[t],
			words:words,
		}
	}

	return tags
}

func generateTagsFromPool(pool []*TagTemplate, min int) []model.KeyValue {
	max := len(pool)
	size := rand.Intn(max-min) + min
	tags := make([]model.KeyValue, size)
	for i := 0; i < size; i++ {
		index := rand.Intn(max - 1)
		tag := pool[index]
		tags[i] = tag.Tag()
	}
	return tags
}
