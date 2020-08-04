package trace

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
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

//client 
func ClientDialOption(parentTracer opentracing.Tracer) grpc.DialOption {
    tracer = parentTracer
    return grpc.WithUnaryInterceptor(grpcClientInterceptor)
}

//text map writer
type TextMapWriter struct {
    metadata.MD
}

//text map writer set
func (t TextMapWriter) Set(key, val string) {
    //key = strings.ToLower(key)
    t.MD[key] = append(t.MD[key], val)
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
    if v := ctx.Value("usergRpcCtx"); v == nil {
    	ctx = ctxShare.Value("usergRpcCtx").(context.Context)
    	logger.Info("jaegerGrpcClientInterceptor ctx", zap.String("parent spanContext", fmt.Sprintf("%s", ctx)))
    }

    var parentContext opentracing.SpanContext
    //从context中获取原始的span
    parentSpan := opentracing.SpanFromContext(ctx)
    if parentSpan != nil {
        parentContext = parentSpan.Context()
    }

    span := tracer.StartSpan(method, opentracing.ChildOf(parentContext))
    defer span.Finish()
    //从context中获取metadata。md.(type) == map[string][]string
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        md = metadata.New(nil)
    } else {
        //如果对metadata进行修改，那么需要用拷贝的副本进行修改。（FromIncomingContext的注释）
        md = md.Copy()
    }
    //定义一个carrier，下面的Inject注入数据需要用到。carrier.(type) == map[string]string
    //carrier := opentracing.TextMapCarrier{}
    carrier := TextMapWriter{md}
    //将span的context信息注入到carrier中
    e := tracer.Inject(span.Context(), opentracing.TextMap, carrier)
    if e != nil {
        logger.Error("tracer Inject err", zap.Error(e))
    }
    //创建一个新的context，把metadata附带上
    ctx = metadata.NewOutgoingContext(ctx, md)
 
    return invoker(ctx, method, req, reply, cc, opts...)
}





