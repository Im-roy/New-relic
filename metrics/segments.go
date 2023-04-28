package metrics

import (
	"context"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type Segment struct{}

func(s Segment) CreateSegment(ctx context.Context, segmentName string) *newrelic.Segment {
	txn := newrelic.FromContext(ctx)
	seg := txn.StartSegment(segmentName)
	return seg
}

func(s Segment) AddAttribute(seg *newrelic.Segment, key string, value interface{}) {
	seg.AddAttribute(key, value)
} 

func(s Segment) AddAttributes(seg *newrelic.Segment, attributes map[string]interface{}) {
	for key, value := range attributes {
		seg.AddAttribute(key, value)
	}
} 

type DatastoreSegment struct{}

func(s DatastoreSegment) CreateDataStoreSegment(ctx context.Context, collection, dbOperationType string) newrelic.DatastoreSegment {
	// input parametres..
	// ctx should have txn, collection is name of table on which operation is being performed.
	// operation is type of operation ex: insert, select etc
	txn := newrelic.FromContext(ctx)
	seg := newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastorePostgres,
		Collection: collection,
		Operation:  dbOperationType,
	}
	return seg
}

func(s DatastoreSegment) AddAttribute(seg *newrelic.Segment, key string, value interface{}) {
	seg.AddAttribute(key, value)
} 

func(s DatastoreSegment) AddAttributes(seg *newrelic.Segment, attributes map[string]interface{}) {
	for key, value := range attributes {
		seg.AddAttribute(key, value)
	}
} 

type ExternalSegment struct{}

func(s ExternalSegment) CreateExternalSegment(ctx context.Context, request *http.Request) *newrelic.ExternalSegment {
	txn := newrelic.FromContext(ctx)
	eseg := newrelic.StartExternalSegment(txn, request)
	return eseg
}

func(s ExternalSegment) AddAttribute(seg *newrelic.Segment, key string, value interface{}) {
	seg.AddAttribute(key, value)
} 

func(s ExternalSegment) AddAttributes(seg *newrelic.Segment, attributes map[string]interface{}) {
	for key, value := range attributes {
		seg.AddAttribute(key, value)
	}
} 

