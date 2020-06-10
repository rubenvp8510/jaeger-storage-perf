package generator

import (
	"github.com/ruben.vp8510/jaeger-storage-perf/generator/data"
	"testing"
)

func TestGenerateTagsKeys(t *testing.T) {
 	/*tags := []string{
		"address",
		"application",
		"category",
		"component",
		"controller",
	}*/
 	keys := generateTagsKeys(data.Tags, 250)
	print(keys)
}
