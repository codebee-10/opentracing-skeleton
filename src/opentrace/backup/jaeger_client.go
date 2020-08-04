package main
import (
	"fmt"
	"io"
	opentracing "github.com/opentracing/opentracing-go"
	// "github.com/opentracing/opentracing-go/ext"
	// metadata "google.golang.org/grpc/metadata"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"golang.org/x/net/context"
	"go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "github.com/natefinch/lumberjack"
	"net/http"
	"io/ioutil"
)

var logger *zap.Logger
var tracer opentracing.Tracer
var closer io.Closer
// sf sampling frequency
var sf = 100

func init() {
	//init log
	initLogger()
	//init tracer
	tracer, closer = InitJaeger("API Client")
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
        Filename:   "logs/opentracing/go/http_trace.log",
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

// SetSamplingFrequency 设置采样频率 0 <= n <= 100
func SetSamplingFrequency(n int) {
	sf = n
}

//log trace
func logTrace(span opentracing.Span) {
	//log trace 
	span.LogKV("event", "awesome trace")
	span.LogKV("func trace event", "log trace")
	
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
	}

	if err := opentracing.GlobalTracer().Inject(
		sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header)); err != nil {
		// log.Println(err)
		fmt.Println("inject error ...", err)
	}
	defer closer.Close()
	defer sp.Finish()
}

//getUserListRequest
func getUserListRequest() string{
	apiUrl := "http://127.0.0.1:10030/getUserList?zone=sz"
	httpClient := &http.Client{}
    r, _ := http.NewRequest("GET", apiUrl, nil)
    //add trace
   	addReqTracer(r)
    response, _ := httpClient.Do(r)
    logger.Info("header....", zap.String("header string", fmt.Sprintf("%s", r.Header)))
    
    if response.StatusCode == 200 {
		str, _ := ioutil.ReadAll(response.Body)
	    bodystr := string(str)
	    logger.Info("response", zap.String("body string", bodystr))
		return bodystr
	 }
	 return "request err..."
}

//run suite
func runSuite() {
	getUserListRequest()
	defer logger.Sync()
}

func main() {
	runSuite()
}





