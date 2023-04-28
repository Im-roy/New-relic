package basics

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"new-relic/metrics"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func MakeExternalUrlCall(req *http.Request) {
	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Process the response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
	log.Println(string(body))
}

func InstrumentExternalSegment(ctx context.Context) {
	txn := newrelic.FromContext(ctx)
	seg := txn.StartSegment("InstrumentExternalSegment")
	defer seg.End()
	URL := "https://bitbucket.org/infracoreplatform/erp-integration-service/commits/c26dac5a98c6936a8dfc291f2c03ce3f430aee9c"
    req, _ := http.NewRequest("GET", URL, nil)
	esegObj := metrics.ExternalSegment{}
	eseg := esegObj.CreateExternalSegment(ctx, req)
	defer eseg.End()
	MakeExternalUrlCall(req)
}

func InstrumentUsingRoundTripper(ctx context.Context) {
	txn := newrelic.FromContext(ctx)
	client := &http.Client{}
	client.Transport = newrelic.NewRoundTripper(client.Transport)

	URL := "https://www.cricbuzz.com/"
	request, _ := http.NewRequest("GET", URL, nil)
	request = newrelic.RequestWithTransactionContext(request, txn)

	resp, err := client.Do(request)
	if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Process the response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
	log.Println(string(body))
}