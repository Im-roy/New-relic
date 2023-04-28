package main

import (
	"context"
	"log"
	"new-relic/basics"
	"testing"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func Test_DataStoreSegmentsWithoutNRPQ(t *testing.T) {
	//*******************************************
	//       Without  nrpq package ...
	// to run without nrpq package uncomment line 85, 105.
	newRelic := NewRelicApp{}
	app := newRelic.StartApp()
	defer app.Shutdown(5 * time.Second)

	db, err := basics.GetDbClient()
	if err != nil {
		panic(err)
	}
	log.Println("db connection established...")

	txn := app.StartTransaction("Postgres transaction... testing without nppq package")
	defer txn.End()
	// create context from txn and pass this to every child methods...
	// - txn := newrelic.FromContext(ctx)
	ctx := newrelic.NewContext(context.Background(), txn)
	// Start a child segment for the database query
	err = basics.InstrumentSelect(ctx, db)
	if err != nil {
		panic(err)
	}
	err = basics.InstrumentInsert(ctx, db)
	if err != nil {
		panic(err)
	}
}

func Test_DataStoreSegmentsWithNRPQ(t *testing.T) {
	//*******************************************
	//      With  nrpq package ...
	//Create data store segments using nrpq package ...
	//- Create a db connection using nrpq driver.
	newRelic := NewRelicApp{}
	app := newRelic.StartApp()
	defer app.Shutdown(5 * time.Second)

	db, err := basics.GetNRGormDbClient()
	if err != nil {
		panic(err)
	}
	log.Println("db connection established...")

	txn := app.StartTransaction("Postgres transaction... testing using nppq package")
	defer txn.End()
	// create context from txn and pass this to every child methods...
	// - txn := newrelic.FromContext(ctx)
	ctx := newrelic.NewContext(context.Background(), txn)
	// Start a child segment for the database query
	err = basics.InstrumentSelectRepoCall(ctx, db)
	if err != nil {
		panic(err)
	}
	err = basics.InstrumentInsertRepoCall(ctx, db)
	if err != nil {
		panic(err)
	}
}

func TestExternalSegmentExample(t *testing.T) {
	newRelic := NewRelicApp{}
	app := newRelic.StartApp()
	defer app.Shutdown(5 * time.Second)

	txn := app.StartTransaction("External transaction... example")
	defer txn.End()
	// create context from txn and pass this to every child methods...
	// - txn := newrelic.FromContext(ctx)
	ctx := newrelic.NewContext(context.Background(), txn)
	basics.InstrumentExternalSegment(ctx)
}

func TestExternalSegmentUsingRoundTripper(t *testing.T) {
	newRelic := NewRelicApp{}
	app := newRelic.StartApp()
	defer app.Shutdown(5 * time.Second)

	txn := app.StartTransaction("External segment using roundtripper... example")
	defer txn.End()
	// create context from txn and pass this to every child methods...
	// - txn := newrelic.FromContext(ctx)
	ctx := newrelic.NewContext(context.Background(), txn)
	basics.InstrumentUsingRoundTripper(ctx)
}

func TestMessageProducerSegement(t *testing.T) {
	newRelic := NewRelicApp{}
	app := newRelic.StartApp()
	defer app.Shutdown(5 * time.Second)

	txn := app.StartTransaction("Message Producer Segement... example")
	defer txn.End()
	// create context from txn and pass this to every child methods...
	// - txn := newrelic.FromContext(ctx)
	ctx := newrelic.NewContext(context.Background(), txn)
	basics.InstrumentMessageProducer(ctx)
}