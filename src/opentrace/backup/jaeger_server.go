package main
import (
	"fmt"
	"io"
	// "time"
	opentracing "github.com/opentracing/opentracing-go"
	// "github.com/opentracing/opentracing-go/ext"
	// metadata "google.golang.org/grpc/metadata"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "github.com/natefinch/lumberjack"
	"golang.org/x/net/context"
	"net/http"
)

var logger *zap.Logger
var tracer opentracing.Tracer
var closer io.Closer

//init
func init() {
	//init log
	initLogger()
	//init tracer
	tracer, closer = InitJaeger("API Gateway")
}

//init logger
func initLogger() {
    writeSyncer := getLogWriter()
    encoder := getEncoder()
    core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
    logger = zap.New(core, zap.AddCaller())
}

//get log writer
func getLogWriter() zapcore.WriteSyncer {
    logMaxSize := 10
    logMaxBackups := 200
    logMaxAge := 30

    lumberJackLogger := &lumberjack.Logger{
        Filename:   "logs/opentracing/go/http_server_trace.log",
        MaxSize:    logMaxSize,
        MaxBackups: logMaxBackups,
        MaxAge:     logMaxAge,
        Compress:   false,
    }
    return zapcore.AddSync(lumberJackLogger)
}

//get log encoder
func getEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    return zapcore.NewConsoleEncoder(encoderConfig)
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
func logTrace(span opentracing.Span) {
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
func baggageTrace(span opentracing.Span) {
	traceStr := "trace awesome thing"
	//use baggage
	// set
	span.SetBaggageItem("greeting", traceStr)
	// get
	greeting := span.BaggageItem("greeting")
	fmt.Println(greeting)
}

//write sub span
func writeSubSpan(span opentracing.Span) {
	//use context
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)

	// 其他过程获取并开始子 span
	newSpan, _ := opentracing.StartSpanFromContext(ctx, "sub span")
	// StartSpanFromContext 会将新span保存到ctx中更新
	defer newSpan.Finish()
}

// TracerWrapper tracer wrapper
func addReqTracer(r *http.Request) {
	opentracing.InitGlobalTracer(tracer)
	sp := tracer.StartSpan(r.URL.Path)
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, 
		opentracing.HTTPHeadersCarrier(r.Header))
	if spanCtx != nil {
		sp = opentracing.GlobalTracer().StartSpan(r.URL.Path, opentracing.ChildOf(spanCtx))
	}else{
		if err := opentracing.GlobalTracer().Inject(
			sp.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header)); err != nil {
			logger.Error("inject error ...", zap.Error(err))
		}
	}

	defer closer.Close()
	defer sp.Finish()
}

//getUserList
func getUserList(w http.ResponseWriter, r *http.Request) {
	logger.Info("server req header", zap.String("request headers", fmt.Sprintf("%s", r.Header)))
	addReqTracer(r)
	io.WriteString(w, "get request")
}

//run suite
func runSuite() {
	fmt.Println("start server ...")
	const HTTP_URL = "0.0.0.0:10030"
	http.HandleFunc("/getUserList", getUserList)

	err := http.ListenAndServe(HTTP_URL, nil)
	if err != nil {
		logger.Error("ListenAndServe err", zap.Error(err))
	}
	fmt.Println("Over ...")
}

//main
func main() {
	runSuite()
}

