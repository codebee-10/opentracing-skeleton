package main
import (
	"fmt"
	"io"
    opentracing "github.com/opentracing/opentracing-go"
    // httpTraceWrapper "github.com/roancsu/traceandtrace-go/libs/trace/wrapper/http"
    httpTraceWrapper "opentracing-skeleton/libs/trace/wrapper/http"
    logger "github.com/roancsu/traceandtrace-go/libs/log"
    "golang.org/x/net/context"
	"net/http"
    "io/ioutil"
)


var tracer opentracing.Tracer
var closer io.Closer
var ctxShare context.Context
var serviceName string


func init() {
    serviceName = "Http Service1"
}


//getUserList
func getUserList(w http.ResponseWriter, r *http.Request) {
    //add trace
    // addHttpTrace(r)
    getUserListRequest(r.Header)
	io.WriteString(w, "get request")
}


//getUserListRequest
func getUserListRequest(header http.Header) string{
    apiUrl := "http://127.0.0.1:10032/getUserList?zone=sz"
    httpClient := &http.Client{}
    r, _ := http.NewRequest("GET", apiUrl, nil)
    r.Header = header
    //add trace
    httpTraceWrapper.AddJaegerTracer(r)

    response, _ := httpClient.Do(r)
    logger.Info(fmt.Sprintf("header string %s", r.Header))
    
    if response.StatusCode == 200 {
        str, _ := ioutil.ReadAll(response.Body)
        bodystr := string(str)
        logger.Info(fmt.Sprintf("body string %s", bodystr))
        return bodystr
     }
     return "request err..."
}



//run suite
func runSuite() {
	fmt.Println("start server ...")
	const HTTP_URL = "0.0.0.0:10031"
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

