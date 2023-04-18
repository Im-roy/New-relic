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
	db, err := basics.GetNRGormDbClient()
	if err != nil {
		log.Fatalf("error in database connection...")
	}
	// db.AutoMigrate(&Product{})
	// ExampleChecker(app, db)
	txn := app.StartTransaction("Postgres ofrm transaction... db")
	defer txn.End()
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
	// End the child segment after the query is executed
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func ExampleChecker(app *newrelic.Application, gormDB *gorm.DB) {
	// Create New Relic Transaction to monitor GORM Database
	gormTransactionTrace := app.StartTransaction("GORM Operation")
	gormTransactionContext := newrelic.NewContext(context.Background(), gormTransactionTrace)
	tracedDB := gormDB.WithContext(gormTransactionContext)

	// Create
	tracedDB.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	tracedDB.First(&product, 1)                 // find product with integer primary key
	tracedDB.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	tracedDB.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	tracedDB.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	tracedDB.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	tracedDB.Delete(&product, 1)

	// End New Relic transaction trace
	gormTransactionTrace.End()
}
