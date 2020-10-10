module opentracing-skeleton

go 1.14

require (
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.4
	github.com/pkg/errors v0.9.1 // indirect
	github.com/roancsu/traceandtrace-go v0.0.0-20200922124606-28115e3f5a3e
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/net v0.0.0-20200904194848-62affa334b73
	google.golang.org/grpc v1.32.0
	src/opentrace/grpc/helloworld v0.0.0-00010101000000-000000000000 // indirect
)

replace src/opentrace/grpc/helloworld => ./src/opentrace/grpc/helloworld/
