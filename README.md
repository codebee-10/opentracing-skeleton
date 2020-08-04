# opentracing-skeleton
opentracing go、java wrapper 


### quick start

1. run grpc
```go
go run src/opentrace/grpc/helloworld/server.go
```

2. run http server
```go
go run src/opentrace/trace_server.go 
```

3. run http client
```go
go run src/opentrace/trace_client.go 
```


### run Jaeger 

```shell
docker run \
-p 5775:5775/udp \
-p 16686:16686 \
-p 6831:6831/udp \
-p 6832:6832/udp \
-p 5778:5778 \
-p 14268:14268 \
jaegertracing/all-in-one:latest
```

### init trace
```go
cfg, err := jaegercfg.FromEnv()  //从环境变量中获取配置信息
cfg.Sampler.Type = "const" 	 //sampler 类型
cfg.Sampler.Param = 1		 //采样速度
cfg.Reporter.LocalAgentHostPort = "127.0.0.1:6831"  //reporter 地址
cfg.Reporter.LogSpans = true  //开启log
```

### spanContext 传递

http

```go
//http 获取span context 信息, 并向下游传递
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
```

grpc

```go
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
//定义span
span := tracer.StartSpan(method, opentracing.ChildOf(parentContext))
defer span.Finish()
```


### reporter 定义



### http.Handler



### rpc.Handler





