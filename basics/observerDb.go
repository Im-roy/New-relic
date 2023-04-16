package basics

import (
	"context"
	"log"

	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDbClient() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
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
		Collection: "playground",
		Operation:  dbOperationType,
	}
	return seg
}

func InstrumentInsert(ctx context.Context, db *gorm.DB) error {
	seg := getDataStoreSegment(ctx, "INSERT")
	defer seg.End()

	err := db.Raw("INSERT INTO playground (type, color, location, install_date) VALUES ('swing', 'blue', 'northwest', '2018-08-16')").Error
	if err != nil {
		return err
	}
	return nil
}

func InstrumentSelect(ctx context.Context, db *gorm.DB) error {
	seg := getDataStoreSegment(ctx, "SELECT")
	defer seg.End()

	var data []int
	err := db.Raw("SELECT equip_id FROM playground").Find(&data).Error
	if err != nil {
		return err
	}
	log.Println(data)
	return nil
}
