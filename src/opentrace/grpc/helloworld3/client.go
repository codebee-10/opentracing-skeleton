package main
 
import (
    "context"
    "log"
    "os"
    "time"
    "io"
    opentracing "github.com/opentracing/opentracing-go"
    jaeger "github.com/uber/jaeger-client-go"
    "github.com/uber/jaeger-client-go/config"
    metadata "google.golang.org/grpc/metadata"
    "google.golang.org/grpc"
    pb "src/opentrace/grpc/helloworld/output/github.com/grpc/example/helloworld"
    "fmt"
)
 
const (
    address     = "localhost:50051"
    defaultName = "ethan"
)


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


func clientDialOption(tracer opentracing.Tracer) grpc.DialOption {
    return grpc.WithUnaryInterceptor(jaegerGrpcClientInterceptor)
}


type TextMapWriter struct {
    metadata.MD
}


func (t TextMapWriter) Set(key, val string) {
    //key = strings.ToLower(key)
    t.MD[key] = append(t.MD[key], val)
}


func jaegerGrpcClientInterceptor (
    ctx context.Context, 
    method string, 
    req, reply interface{},
    cc *grpc.ClientConn, 
    invoker grpc.UnaryInvoker, 
    opts ...grpc.CallOption) (err error) {

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

 
func main() {
    tracer, closer, err := initJaeger("rpc client", "127.0.0.1:6831")
    if err != nil {
        // log.Fatal(err)
        fmt.Println("init jaeger err", err)
    }
    defer closer.Close()
    // Set up a connection to the server.
    conn, err := grpc.Dial(address, grpc.WithInsecure(), clientDialOption(tracer))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewGreeterClient(conn)
 
    // Contact the server and print out its response.
    name := defaultName
    if len(os.Args) > 1 {
        name = os.Args[1]
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("Greeting: %s", r.Message)
}
















