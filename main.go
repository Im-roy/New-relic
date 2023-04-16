package main

import (
	"context"
	"fmt"
	"log"
	"new-relic/basics"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func makeApplication(name string) (*newrelic.Application, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(name),
		newrelic.ConfigLicense("38d5a8e5ab161d51bfdc60bf00834d810971NRAL"),
	)
	if nil != err {
		return nil, err
	}

	// Wait for the application to connect.
	if err = app.WaitForConnection(10 * time.Second); nil != err {
		return nil, err
	}

	return app, nil
}

func main() {
	app, err := makeApplication("Example AppV3")
	if err != nil {
		panic(err)
	}
	fmt.Println("new relic connection established....")
	defer app.Shutdown(5 * time.Second)

	// Create a sample transaction...
	// basics.Print(app)

	// Create a transaction and segment both...
	// msg, err := basics.Hello(app, "Akash Yadav")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(msg)

	//*******************************
	//  Create postgres db connection..
	db, err := basics.GetDbClient()
	if err != nil {
		log.Fatalf("error in database connection...")
	}
	txn := app.StartTransaction("Postgres transaction...")
	defer txn.End()
	ctx := newrelic.NewContext(context.Background(), txn)
	// Start a child segment for the database query
	err = basics.InstrumentInsert(ctx, db)
	if err != nil {
		panic(err)
	}
	// err = basics.InstrumentSelect(ctx, db)
	// if err != nil {
	// 	panic(err)
	// }
	// End the child segment after the query is executed
}
