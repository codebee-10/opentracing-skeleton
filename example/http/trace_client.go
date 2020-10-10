package main
import (
	"fmt"
	logger "github.com/roancsu/traceandtrace-go/libs/log"
	// httpTraceWrapper "github.com/roancsu/traceandtrace-go/libs/trace/wrapper/http"
	httpTraceWrapper "opentracing-skeleton/libs/trace/wrapper/http"
	"net/http"
	"io/ioutil"
)

 
//getUserListRequest
func getUserListRequest() string{
	apiUrl := "http://127.0.0.1:10031/getUserList?zone=sz"
	httpClient := &http.Client{}
    r, _ := http.NewRequest("GET", apiUrl, nil)

    //add trace
    //jaeger
	httpTraceWrapper.AddJaegerTracer(r, "UserList")
	//zipkin
	// httpTraceWrapper.AddZipkinTracer(r, "API Client2")

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
	getUserListRequest()
}


func main() {
	runSuite()
}





