package trace

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
    "golang.org/x/net/context"
    zaplog "libs/log"
    "go.uber.org/zap"
    "io"
    "net/http"
    "fmt"
)


var tracer opentracing.Tracer
// var closer io.Closer
var ctxShare context.Context
var logger *zap.Logger

var sf = 100

const (
    address     = "localhost:50051"
    defaultName = "ethan"
)

//init
func init() {
    //init log
    logger = zaplog.InitLogger()
}

// init Jaeger
func InitJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg, err := jaegercfg.FromEnv()
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1
	cfg.Reporter.LocalAgentHostPort = "127.0.0.1:6831"
	cfg.Reporter.LogSpans = true
	
	tracer, closer, err := cfg.New(service, jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

//log trace
func LogTrace(span opentracing.Span) {
	//log trace 
	span.LogKV("event", "git service server")
	//log trace
	span.SetTag("value", "traceTag")
	span.LogFields(
	    log.String("event", "awesome report"),
	    log.String("value", "traceTag"),
	)
}

//baggage trace
func BaggageTrace(span opentracing.Span) {
	traceStr := "trace awesome thing"
	//use baggage
	// set
	span.SetBaggageItem("greeting", traceStr)
	// get
	greeting := span.BaggageItem("greeting")
	fmt.Println(greeting)
}

//write sub span
func WriteSubSpan(span opentracing.Span) {
	//use context
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)

	// 其他过程获取并开始子 span
	newSpan, _ := opentracing.StartSpanFromContext(ctx, "sub span")
	// StartSpanFromContext 会将新span保存到ctx中更新
	defer newSpan.Finish()
}

// TracerWrapper tracer wrapper
func AddReqTracer(r *http.Request, tracer opentracing.Tracer) {
	opentracing.InitGlobalTracer(tracer)
	sp := tracer.StartSpan(r.URL.Path)
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, 
		opentracing.HTTPHeadersCarrier(r.Header))
	if spanCtx != nil {
		sp = opentracing.GlobalTracer().StartSpan(r.URL.Path, opentracing.ChildOf(spanCtx))
	}else{
		//http inject
		if err := opentracing.GlobalTracer().Inject(
			sp.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header)); err != nil {
			logger.Error("inject error ...", zap.Error(err))
		}
	}

	//上下文记录父spanContext
	ctxShare = context.WithValue(context.Background(), "usergRpcCtx", 
        opentracing.ContextWithSpan(context.Background(), sp))

	defer sp.Finish()
}

















