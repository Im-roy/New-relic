package basics

import (
	"context"
	"log"

	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetNRGormDbClient() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=dev port=5432 sslmode=disable"

	// Create a dialector object with Driver name as nrpostgres and connection string.
	dialector := postgres.Dialector{
		Config: &postgres.Config{
			DriverName: "nrpostgres",
			DSN:        dsn,
		},
	}
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return gormDB, nil
}

func InstrumentSelectRepoCall(ctx context.Context, db *gorm.DB) error {
	//******* Create a data store segment using nrpq...
	// step-1: create a updated gorm.Db client 
	// Name of span is InstrumentSelectRepoCall
	// Name of data store segment = Db-type + table_name + type of DDL/DML command.
	tracedDB := db.WithContext(ctx) // step-1
	// extracting txn from context and creating a segment to observe this method.
	// txn := newrelic.FromContext(ctx)
	// tseg := txn.StartSegment("InstrumentSelectRepoCall") 
	// defer tseg.End()
	var data []int
	err := tracedDB.Table("erp.request_response_audits").Select("Id").Find(&data).Error
	if err != nil {
		return err
	}
	log.Println(data)
	return nil
}

func InstrumentInsertRepoCall(ctx context.Context, db *gorm.DB) error {

	//******* Create a data store segment using nrpq...
	// step-1: create a updated gorm.Db client 
	// Name of span is InstrumentSelectRepoCall
	// Name of data store segment = Db-type + table_name + type of DDL/DML command.

	tracedDB := db.WithContext(ctx) // step-1
	// extracting txn from context and creating a segment to observe this method.
	// txn := newrelic.FromContext(ctx)
	// tseg := txn.StartSegment("InstrumentInsertRepoCall")
	// defer tseg.End()
	err := tracedDB.Debug().Table("erp.request_response_audits").Create(map[string]interface{}{
		"transaction_id": "134567", 
		"request_id": "jifjdfnufjfnjf",
		"source_service": "test_service",
		"event_type": "test-event-type",
		"event_sub_type": "test-event-sub-type",
	}).Error
	if err != nil {
		panic(err)
	}
	return nil
}
