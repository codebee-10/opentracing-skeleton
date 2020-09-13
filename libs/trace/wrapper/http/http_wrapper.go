package trace

import (
	opentracing "github.com/opentracing/opentracing-go"
	jaegertrace "libs/trace/reporter/jaeger"
    zipkintrace "libs/trace/reporter/zipkin"
    "net/http"
)


//jaeger
func AddJaegerTracer(r *http.Request, serviceName string) opentracing.Tracer{
    tracer, closer := jaegertrace.InitJaeger(serviceName)
    defer closer.Close()
    jaegertrace.AddTracer(r, tracer)
    return tracer
}

//zipkin
func AddZipkinTracer(r *http.Request, serviceName string) opentracing.Tracer{
    //TO-DO
    tracer, closer := zipkintrace.InitZipkin(serviceName)
    defer closer.Close()
    zipkintrace.AddTracer(r, tracer)
    return tracer
}

//skywalking
func AddSkyWalkingTracer(r *http.Request, serviceName string) {
    //TO-DO
}
