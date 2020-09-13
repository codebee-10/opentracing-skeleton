package main

import (
	"log"
	zipkin "github.com/openzipkin/zipkin-go"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"fmt"
)


func main() {
	// set up a span reporter
	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	fmt.Println(reporter)
	defer reporter.Close()

	// create endpoint
	endpoint, err := zipkin.NewEndpoint("API Client", "")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// create tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	sp := tracer.StartSpan("test service")

	defer sp.Finish()

}