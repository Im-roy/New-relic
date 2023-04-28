package basics

import (
	"context"
	"log"
	
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDbClient() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=dev port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}
	return db, err
}

func getDataStoreSegment(ctx context.Context, dbOperationType string) newrelic.DatastoreSegment {
	txn := newrelic.FromContext(ctx)
	seg := newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastorePostgres,
		Collection: "erp.request_response_audits",
		Operation:  dbOperationType,
	}
	return seg
}

func InstrumentInsert(ctx context.Context, db *gorm.DB) error {
	seg := getDataStoreSegment(ctx, "INSERT")
	defer seg.End()

	err := db.Debug().Table("erp.request_response_audits").Create(map[string]interface{}{
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

func InstrumentSelect(ctx context.Context, db *gorm.DB) error {
	seg := getDataStoreSegment(ctx, "SELECT")
	defer seg.End()

	var data []int
	err := db.Table("erp.request_response_audits").Select("Id").Find(&data).Error
	if err != nil {
		return err
	}
	log.Println(data)
	return nil
}
