module go-worker-pools

go 1.14

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/charithe/otgrpc v0.0.0-20170514181245-1f3477f51faf // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/opentracing-contrib/go-grpc v0.0.0-20191001143057-db30781987df // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5 // indirect
	github.com/spf13/viper v1.7.0 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.uber.org/zap v1.15.0 // indirect
	google.golang.org/genproto v0.0.0-20200731012542-8145dea6a485 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	libs/log v0.0.0-00010101000000-000000000000 // indirect
	libs/trace v0.0.0-00010101000000-000000000000 // indirect
	src/opentrace/grpc/helloworld v0.0.0-00010101000000-000000000000 // indirect
)

replace src/opentrace/grpc/helloworld => ./src/opentrace/grpc/helloworld/

replace libs/log => ./libs/log

replace libs/trace => ./libs/trace
