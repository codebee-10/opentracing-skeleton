package trace

import (
	"google.golang.org/grpc"
    jaegertrace "opentracing-skeleton/libs/trace/reporter/jaeger"
    "net/http"
    "io"
)


//jaeger
func AddJaegerTracer(serviceName string) (grpc.DialOption, io.Closer) {
    tracer, closer := jaegertrace.InitJaeger(serviceName)
    return jaegertrace.ClientDialOption(tracer), closer
}


//zipkin
func AddZipkinTracer(r *http.Request, serviceName string) {
    //TO-DO
}


//skywalking
func AddSkyWalkingTracer(r *http.Request, serviceName string) {
    //TO-DO
}

 