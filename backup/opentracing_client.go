func main() {
	//init jaeger
	tracer, closer, err := initJaeger("client", jaegerAgentHost)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()
	//dial
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), clientDialOption(tracer))
	if err != nil {
		log.Fatalf("dial fail, %+v\n", err)
	}
	//发送请求
	req := &delayqueue.PingRequest{Msg:"ping~"}
	client := delayqueue.NewDelayQueueClient(conn)
	r, err := client.Ping(context.Background(), req)
 
	fmt.Println(r, err)
}
 
func clientDialOption(tracer opentracing.Tracer) grpc.DialOption {
	return grpc.WithUnaryInterceptor(jaegerGrpcClientInterceptor)
}


type TextMapWriter struct {
	metadata.MD
}
//重写TextMapWriter的Set方法，我们需要将carrier中的数据写入到metadata中，这样grpc才会携带。
func (t TextMapWriter) Set(key, val string) {
	//key = strings.ToLower(key)
	t.MD[key] = append(t.MD[key], val)
}
 
func jaegerGrpcClientInterceptor (ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	var parentContext opentracing.SpanContext
	//先从context中获取原始的span
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		parentContext = parentSpan.Context()
	}
	tracer := opentracing.GlobalTracer()
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
		fmt.Println("tracer Inject err,", e)
	}
	//创建一个新的context，把metadata附带上
	ctx = metadata.NewOutgoingContext(ctx, md)
 
	return invoker(ctx, method, req, reply, cc, opts...)
}
 
func initJaeger(service string, jaegerAgentHost string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort:jaegerAgentHost,
		},
	}
	tracer, closer, err = cfg.New(service, config.Logger(jaeger.StdLogger))
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, err
}
