package zipkin

import (
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go/reporter"
	"golang.org/x/net/context"
	"go.uber.org/zap"
	zaplog "libs/log"
	"net/http"
	"log"
)

var logger *zap.Logger
var ctxShare context.Context
var tracer opentracing.Tracer

//init
func init() {
	 logger = zaplog.InitLogger()
}

//init zipkin
func InitZipkin(service string) (opentracing.Tracer, reporter.Reporter){
	// set up a span reporter
	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(service, "")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	tracer := zipkinot.Wrap(nativeTracer)
	return tracer, reporter
}

//add tracer
func AddTracer(r *http.Request, tracer opentracing.Tracer) {
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
	ctxShare = context.WithValue(context.Background(), "usergRpcCtx", opentracing.ContextWithSpan(context.Background(), sp))

	defer sp.Finish()
}





