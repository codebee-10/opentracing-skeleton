package trace

import (
	"google.golang.org/grpc"
    jaegertrace "libs/trace/reporter/jaeger"
    "net/http"
)


//jaeger
func AddJaegerTracer(serviceName string) grpc.DialOption {
    tracer, _ := jaegertrace.InitJaeger(serviceName)
    // defer closer.Close()
    return jaegertrace.ClientDialOption(tracer)
}

//zipkin
func AddZipkinTracer(r *http.Request, serviceName string) {
    //TO-DO
}

//skywalking
func AddSkyWalkingTracer(r *http.Request, serviceName string) {
    //TO-DO
}
