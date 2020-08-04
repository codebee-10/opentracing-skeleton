package main
import (
	"fmt"
	"io"
    opentracing "github.com/opentracing/opentracing-go"
    jaegertrace "libs/trace/jaeger"
    "go.uber.org/zap"
    zaplog "libs/log"
    pb "src/opentrace/grpc/helloworld/output/github.com/grpc/example/helloworld"
    "google.golang.org/grpc"
    "golang.org/x/net/context"
	"net/http"
    "os"
    "time"
)


var logger *zap.Logger
var tracer opentracing.Tracer
var closer io.Closer
var ctxShare context.Context


const (
    address     = "localhost:50051"
    defaultName = "ethan"
)


//init
func init() {
	//init log
	logger = zaplog.InitLogger()
    //init tracer
    tracer, closer = jaegertrace.InitJaeger("API Gateway")
}

//grpc request
func getUserListRpcRequest(tracer opentracing.Tracer) {
    conn, err := grpc.Dial(address, grpc.WithInsecure(),jaegertrace.ClientDialOption(tracer))
    if err != nil {
        logger.Error("did not connect", zap.Error(err))
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
        logger.Error("could not greet", zap.Error(err))
    }
    fmt.Println("Greeting: %s", r.Message)
    logger.Info("Greeting", zap.String("message", r.Message))
}

//getUserList
func getUserList(w http.ResponseWriter, r *http.Request) {
	logger.Info("server req header", zap.String("request headers", fmt.Sprintf("%s", r.Header)))
	//http request
	jaegertrace.AddReqTracer(r, tracer)
	logger.Info("httpRequest ....")

	//get user list rpc request
	getUserListRpcRequest(tracer)

	logger.Info("grpcRequest ....")
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

