package main

import(
    pb "src/opentrace/grpc/helloworld/output/github.com/grpc/example/helloworld"
    "google.golang.org/grpc"
    logger "github.com/roancsu/traceandtrace-go/libs/log"
    rpcTraceWrapper "github.com/roancsu/traceandtrace-go/libs/trace/wrapper/rpc"
    "golang.org/x/net/context"
    "fmt"
    "time"
)


const (
    addr     = "localhost:50050"
    serviceName = "rpc client"
)



func main() {
    opt, closer := rpcTraceWrapper.AddJaegerTracer(serviceName)
    defer closer.Close()
    // dial
    conn, err := grpc.Dial(addr, grpc.WithInsecure(), opt)
    if err != nil {
    }
    //发送请求
    name := "ethan"
    ctx, _ := context.WithTimeout(context.Background(), time.Second)
    c := pb.NewGreeterClient(conn)
    r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
    if err != nil {
        logger.Error(fmt.Sprintf("could not greet %s", err))
    }
    fmt.Println("Greeting: %s", r.Message)
}




