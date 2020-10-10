package main
import (
	"fmt"
	"io"
    opentracing "github.com/opentracing/opentracing-go"
    // httpTraceWrapper "github.com/roancsu/traceandtrace-go/libs/trace/wrapper/http"
    // rpcTraceWrapper "github.com/roancsu/traceandtrace-go/libs/trace/wrapper/rpc"
    httpTraceWrapper "opentracing-skeleton/libs/trace/wrapper/http"
    rpcTraceWrapper "opentracing-skeleton/libs/trace/wrapper/rpc"
    pb "src/opentrace/grpc/helloworld/output/github.com/grpc/example/helloworld"
    "google.golang.org/grpc"
    "golang.org/x/net/context"
	"net/http"
    "os"
    "time"
)


var tracer opentracing.Tracer
var closer io.Closer
var ctxShare context.Context
var serviceName string


const (
    address     = "localhost:50050"
    defaultName = "ethan"
)


func init() {
    serviceName = "Http Service3"
}


//addHttpTrace
func addHttpTrace(r *http.Request) (opentracing.Tracer){
    return httpTraceWrapper.AddJaegerTracer(r, serviceName)
}


func addRpcTrace() grpc.DialOption{
    opt, _ := rpcTraceWrapper.AddJaegerTracer(serviceName)
    return opt
}


//grpc request
func getUserListRpcRequest(tracerOption grpc.DialOption) {
    conn, err := grpc.Dial(address, grpc.WithInsecure(), tracerOption)
    if err != nil {
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
    }
    fmt.Println("Greeting: %s", r.Message)
}


//getUserList
func getUserList(w http.ResponseWriter, r *http.Request) {
	// logger.Info("server req header", zap.String("request headers", fmt.Sprintf("%s", r.Header)))
    //add trace
    addHttpTrace(r)
	//user list rpc request
    getUserListRpcRequest(addRpcTrace()) 
	// logger.Info("grpcRequest ....")
	io.WriteString(w, "get request")
}


//run suite
func runSuite() {
	fmt.Println("start server ...")
	const HTTP_URL = "0.0.0.0:10033"
	http.HandleFunc("/getUserList", getUserList)
	err := http.ListenAndServe(HTTP_URL, nil)
	if err != nil {
		// logger.Error("ListenAndServe err", zap.Error(err))
	}
	fmt.Println("Over ...")
}


//main
func main() {
	runSuite()
}

