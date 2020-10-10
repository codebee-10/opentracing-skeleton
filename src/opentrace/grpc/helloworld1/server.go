package main
 
import (
    "context"
    "log"
    "net"
    "io"
    opentracing "github.com/opentracing/opentracing-go"
    jaeger "github.com/uber/jaeger-client-go"
    "github.com/uber/jaeger-client-go/config"
    metadata "google.golang.org/grpc/metadata"
    "google.golang.org/grpc"
    pb "src/opentrace/grpc/helloworld/output/github.com/grpc/example/helloworld"
    "google.golang.org/grpc/reflection"
    "fmt"
)
 
const (
    port = ":50051"
    serviceName = "rpc:server1:client"
)
 
// server is used to implement helloworld.GreeterServer.
type server struct{}
 
// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{Message: "Hello " + in.Name}, nil
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



func serverOption(tracer opentracing.Tracer) grpc.ServerOption {
    return grpc.UnaryInterceptor(jaegerGrpcServerInterceptor)
}

 
type TextMapReader struct {
    metadata.MD
}


//读取metadata中的span信息
func (t TextMapReader) ForeachKey(handler func(key, val string) error) error { //不能是指针
    for key, val := range t.MD {
        for _, v := range val {
            if err := handler(key, v); err != nil {
                return err
            }
        }
    }
    return nil
}


func jaegerGrpcServerInterceptor(
    ctx context.Context, 
    req interface{}, 
    info *grpc.UnaryServerInfo, 
    handler grpc.UnaryHandler) (resp interface{}, err error) {
    //从context中获取metadata。md.(type) == map[string][]string
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        md = metadata.New(nil)
    } else {
        //如果对metadata进行修改，那么需要用拷贝的副本进行修改。（FromIncomingContext的注释）
        md = md.Copy()
    }
    carrier := TextMapReader{md}
    tracer := opentracing.GlobalTracer()
    spanContext, e := tracer.Extract(opentracing.TextMap, carrier)
    if e != nil {
        fmt.Println("Extract err:", e)
    }
 
    span := tracer.StartSpan(info.FullMethod, opentracing.ChildOf(spanContext))
    defer span.Finish()
    fmt.Println(span)
    ctx = opentracing.ContextWithSpan(ctx, span)
 
    return handler(ctx, req)
}


func main() {
    fmt.Println("rpc server start ...")
    tracer, closer, err := initJaeger(serviceName, "127.0.0.1:6831")
    if err != nil {
        // log.Fatal(err)
        fmt.Println("init jaeger err", err)
    }
    defer closer.Close()

    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    opts := serverOption(tracer)
    s := grpc.NewServer(opts)
    pb.RegisterGreeterServer(s, &server{})
    // Register reflection service on gRPC server.
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
 