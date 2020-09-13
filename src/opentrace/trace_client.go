package main
import (
	"fmt"
	"go.uber.org/zap"
	zaplog "libs/log"	
	"net/http"
	"io/ioutil"
	// opentracing "github.com/opentracing/opentracing-go"
	httpTraceWrapper "libs/trace/wrapper/http"
)

var logger *zap.Logger
// var tracer opentracing.Tracer
// var closer io.Closer

//init
func init() {
	//init log
	logger = zaplog.InitLogger()
}

//addHttpTrace
func addHttpTrace(r *http.Request) {
	// httpTraceWrapper.AddJaegerTracer(r, "API Client")
	httpTraceWrapper.AddZipkinTracer(r, "API Client2")
}

//getUserListRequest
func getUserListRequest() string{
	apiUrl := "http://127.0.0.1:10030/getUserList?zone=sz"
	httpClient := &http.Client{}
    r, _ := http.NewRequest("GET", apiUrl, nil)
    //add trace
    addHttpTrace(r)

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





