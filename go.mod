module github.com/rubenvp8510/jaeger-storage-perf

go 1.13

replace github.com/rubenvp8510/redbull => /home/rvargasp/go/src/github.com/rubenvp8510/redbull

replace github.com/jaegertracing/jaeger v1.18.1 => /home/rvargasp/go/src/github.com/rubenvp8510/jaeger

require (
	github.com/google/uuid v1.1.1
	github.com/jaegertracing/jaeger v1.18.1
	github.com/jaegertracing/jaeger-idl v0.0.0-20200626175211-52fb4c944067 // indirect
	github.com/rubenvp8510/jaeger-storages v0.0.0-20200706195738-7040b4c75869
	github.com/rubenvp8510/redbull v0.0.0-20200703030353-bca447f62bf0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	github.com/uber/jaeger-lib v2.2.0+incompatible
	go.uber.org/zap v1.16.0
)
