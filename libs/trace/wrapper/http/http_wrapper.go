package trace

import (
	opentracing "github.com/opentracing/opentracing-go"
	jaegertrace "opentracing-skeleton/libs/trace/reporter/jaeger"
    zipkintrace "opentracing-skeleton/libs/trace/reporter/zipkin"
    logger "github.com/roancsu/traceandtrace-go/libs/log"
    "strings"
    "net/http"
)


//jaeger
func AddJaegerTracer(r *http.Request, serviceName ...string) opentracing.Tracer{
    //获取服务名称
    svcName := r.URL.Path
    if len(serviceName) != 0 {
        svcName = serviceName[len(serviceName)-1]
    }
    svcName = strings.Replace(svcName, "/", "", -1)
    logger.Info(svcName)
    //初始化jaeger
    tracer, closer := jaegertrace.InitJaeger(svcName)
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
