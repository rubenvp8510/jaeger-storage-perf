module github.com/rubenvp8510/jaeger-storage-perf

go 1.13

replace github.com/rubenvp8510/jaeger-storages => /home/rvargasp/go/src/github.com/rubenvp8510/jaeger-storages

replace github.com/rubenvp8510/redbull => /home/rvargasp/go/src/github.com/rubenvp8510/redbull

require (
	github.com/golang/protobuf v1.3.4
	github.com/google/uuid v1.1.1
	github.com/jaegertracing/jaeger v1.18.1
	github.com/rubenvp8510/jaeger-storages v0.0.0-00010101000000-000000000000
	github.com/rubenvp8510/redbull v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v1.7.0
	github.com/uber/jaeger-lib v2.2.0+incompatible
	go.uber.org/atomic v1.6.0
	go.uber.org/zap v1.15.0
)
