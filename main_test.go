package main

import (
	"flag"
	"github.com/jaegertracing/jaeger/pkg/config"
	"github.com/jaegertracing/jaeger/plugin/storage/badger"
	"github.com/jaegertracing/jaeger/storage"
	"github.com/rubenvp8510/jaeger-storage-perf/helpers"
	"github.com/rubenvp8510/jaeger-storages/druid"
	"github.com/rubenvp8510/jaeger-storages/questbd"
	"github.com/rubenvp8510/redbull/pkg/redbull"
	"github.com/spf13/viper"
	"os"
	"testing"
)

var storageType = flag.String("storage-type", "badger", "storage-type")

var factories map[string]helpers.SetupFunction

func BenchmarkStorage(t *testing.B) {
	stype := *storageType
	factory, ok := factories[stype]
	if !ok {
		panic("Cannot construct factory of type " + stype)
	}
	helpers.StorageBenchmark(stype, factory(), t)
}

func TestMain(m *testing.M) {
	factories = map[string]helpers.SetupFunction{
		"badger": func() storage.Factory {
			factory := badger.NewFactory()
			opts := badger.NewOptions("badger")
			v, _ := config.Viperize(opts.AddFlags)
			factory.InitFromViper(v)
			return factory
		},
		"redbull": func() storage.Factory {
			factory := redbull.NewFactory()
			vip := viper.New()
			vip.Set("redbull.sybil-path", "/home/rvargasp/go/bin/sybil")
			factory.InitFromViper(vip)
			return factory
		},
		"questdb": func() storage.Factory {
			return questbd.NewFactory()
		},
		"druid": func() storage.Factory {
			factory := druid.NewFactory()
			factory.InitFromOptions(druid.DefaultOptions())
			return factory
		},
	}
	helpers.Setup("traces_fixtures")
	code := m.Run()
	os.Exit(code)
}
