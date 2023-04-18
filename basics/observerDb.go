package basics

import (
	"context"
	"log"
	"time"

	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
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

func GetNRGormDbClient() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
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

func InstrumentSelectRepoCall(ctx context.Context, db *gorm.DB) error {
	tracedDB := db.WithContext(ctx)
	txn := newrelic.FromContext(ctx)
	tseg := txn.StartSegment("InstrumentSelectRepoCall")
	defer tseg.End()
	var data []int
	err := tracedDB.Raw("SELECT equip_id FROM playground").Find(&data).Error
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	log.Println(data)
	return nil
}

func InstrumentInsertRepoCall(ctx context.Context, db *gorm.DB) error {
	tracedDB := db.WithContext(ctx)
	txn := newrelic.FromContext(ctx)
	tseg := txn.StartSegment("InstrumentInsertRepoCall")
	defer tseg.End()
	err := tracedDB.Raw("INSERT INTO playground (type, color, location, install_date) VALUES ('swing', 'browm', 'northwest', '2018-08-16')").Error
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	return nil
}
