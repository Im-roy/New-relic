package main

import (
	"context"
	"fmt"
	"log"
	"new-relic/basics"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
)

type NewRelicApp struct {
	App *newrelic.Application
}

func(n *NewRelicApp)StartApp() *newrelic.Application {
	app, err := makeApplication("Example AppV4")
	if err != nil {
		panic(err)
	}
	n.App = app
	return n.App
}

func makeApplication(name string) (*newrelic.Application, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(name),
		newrelic.ConfigLicense("38d5a8e5ab161d51bfdc60bf00834d810971NRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
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
	app, err := makeApplication("Example AppV4")
	if err != nil {
		panic(err)
	}
	fmt.Println("new relic connection established....")
	defer app.Shutdown(5 * time.Second)

	//*******************************
	// Create a sample transaction...
	// basics.Print(app)

	// Create a transaction and segment both...
	// msg, err := basics.Hello(app, "Akash Yadav")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(msg)

	//*******************************************
	//       With  nrpq package ...
	// Create data store segments using nrpq package ...
	// - Create a db connection using nrpq driver.
	// uncomment line - 59, 80.

	// db, err := basics.GetNRGormDbClient()
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("db connection established...")

	// txn := app.StartTransaction("Postgres transaction... testing using nppq package")
	// defer txn.End()
	// // create context from txn and pass this to every child methods...
	// // - txn := newrelic.FromContext(ctx)
	// ctx := newrelic.NewContext(context.Background(), txn)
	// // Start a child segment for the database query
	// err = basics.InstrumentSelectRepoCall(ctx, db)
	// if err != nil {
	// 	panic(err)
	// }
	// err = basics.InstrumentInsertRepoCall(ctx, db)
	// if err != nil {
	// 	panic(err)
	// }
	//*******************************************
	//       Without  nrpq package ...
	// to run without nrpq package uncomment line 85, 105.

	db, err := basics.GetDbClient()
	if err != nil {
		panic(err)
	}
	log.Println("db connection established...")

	txn := app.StartTransaction("Postgres transaction... testing without nppq package", nil, nil)
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

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// func ExampleChecker(app *newrelic.Application, gormDB *gorm.DB) {
// 	// Create New Relic Transaction to monitor GORM Database
// 	gormTransactionTrace := app.StartTransaction("GORM Operation")
// 	gormTransactionContext := newrelic.NewContext(context.Background(), gormTransactionTrace)
// 	tracedDB := gormDB.WithContext(gormTransactionContext)

// 	// Create
// 	tracedDB.Create(&Product{Code: "D42", Price: 100})

// 	// Read
// 	var product Product
// 	tracedDB.First(&product, 1)                 // find product with integer primary key
// 	tracedDB.First(&product, "code = ?", "D42") // find product with code D42

// 	// Update - update product's price to 200
// 	tracedDB.Model(&product).Update("Price", 200)
// 	// Update - update multiple fields
// 	tracedDB.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
// 	tracedDB.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

// 	// Delete - delete product
// 	tracedDB.Delete(&product, 1)

// 	// End New Relic transaction trace
// 	gormTransactionTrace.End()
// }
