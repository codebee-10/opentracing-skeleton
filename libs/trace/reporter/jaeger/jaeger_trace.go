package trace

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "golang.org/x/net/context"
    logger "github.com/roancsu/traceandtrace-go/libs/log"
    "os"
    "io"
    "net/http"
    "fmt"
)


var tracer opentracing.Tracer
var ctxShare context.Context
var rpcCtx string
var sf = 100


//初始化 Jaeger
func InitJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg, err := jaegercfg.FromEnv()
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1
	cfg.Reporter.LocalAgentHostPort = "127.0.0.1:6831"
	if agentHost := os.Getenv("TRACE_AGENT_HOST"); agentHost!="" {
		cfg.Reporter.LocalAgentHostPort = agentHost
	}
	cfg.Reporter.LogSpans = true
	tracer, closer, err := cfg.New(service, jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		logger.Error(fmt.Sprintf("cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
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
	//set
	span.SetBaggageItem("greeting", traceStr)
	//get
	greeting := span.BaggageItem("greeting")
	logger.Info(fmt.Sprintf("greeting: %v\n", greeting))
}


//write sub span
func WriteSubSpan(span opentracing.Span, subSpanName string) {
	//use context
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	//其他过程获取并开始子 span
	newSpan, _ := opentracing.StartSpanFromContext(ctx, subSpanName)
	//StartSpanFromContext 会将新span保存到ctx中更新
	defer newSpan.Finish()
}


// TracerWrapper tracer wrapper
func AddTracer(r *http.Request, tracer opentracing.Tracer) {
	opentracing.InitGlobalTracer(tracer)
	var sp opentracing.Span
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, 
		opentracing.HTTPHeadersCarrier(r.Header))
	if spanCtx != nil {
		sp = opentracing.GlobalTracer().StartSpan(r.URL.Path, opentracing.ChildOf(spanCtx))
		logger.Error(fmt.Sprintf("parent ....: %v\n", r.URL.Path))
	}else{
		sp = tracer.StartSpan(r.URL.Path)
		logger.Error(fmt.Sprintf("new ....: %v\n", r.URL.Path))
	}
	//注入span
	if err := opentracing.GlobalTracer().Inject(
		sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header)); err != nil {
        logger.Error(fmt.Sprintf("inject failed ...: %v\n", err))
	}

	//上下文记录父spanContext
	rpcCtx = r.URL.Path
	ctxShare = context.WithValue(context.Background(), rpcCtx, opentracing.ContextWithSpan(context.Background(), sp))
    //close span
	defer sp.Finish()
}


//client 
func ClientDialOption(parentTracer opentracing.Tracer) grpc.DialOption {
    tracer = parentTracer
    return grpc.WithUnaryInterceptor(grpcClientInterceptor)
}


//grpcClientInterceptor
func grpcClientInterceptor (
    ctx context.Context, 
    method string, 
    req, reply interface{},
    cc *grpc.ClientConn, 
    invoker grpc.UnaryInvoker, 
    opts ...grpc.CallOption) (err error) {

    //上下文获取spanContext
    if rpcCtx != "" {
    	if v := ctx.Value(rpcCtx); v == nil {
	    	ctx = ctxShare.Value(rpcCtx).(context.Context)
	        logger.Info(fmt.Sprintf("trace rpc parent ctx ... %v\n", ctx))
	    }
    }

    //从context中获取metadata
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        md = metadata.New(nil)
    } else {
        //如果对metadata进行修改，那么需要用拷贝的副本进行修改
        md = md.Copy()
    }
    //carrier := opentracing.TextMapCarrier{}
    carrier := TextMapWriter{md}
    //父类 context
    var currentContext opentracing.SpanContext
    //从context中获取原始的span
    parentSpan := opentracing.SpanFromContext(ctx)
    if parentSpan != nil {
        currentContext = parentSpan.Context()
    }else{
    	//start span
	    span := tracer.StartSpan(method)
	    defer span.Finish()
	    currentContext = span.Context()
    }
    
    //将span的context信息注入到carrier中
    e := tracer.Inject(currentContext, opentracing.TextMap, carrier)
    if e != nil {
        logger.Error(fmt.Sprintf("tracer inject failed ...: %v\n", e))
    }
    //创建一个新的context，把metadata附带上
    ctx = metadata.NewOutgoingContext(ctx, md)
    return invoker(ctx, method, req, reply, cc, opts...)
}


//text map writer
type TextMapWriter struct {
    metadata.MD
}


//text map writer set
func (t TextMapWriter) Set(key, val string) {
    t.MD[key] = append(t.MD[key], val)
}






