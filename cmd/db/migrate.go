package main

import (
	"context"
	"flag"
	"go.uber.org/fx"
	"log"
	"uroborus/common/logging"
	settings "uroborus/common/setting"
	"uroborus/model"
	"uroborus/store"
	storeFx "uroborus/store/fx"
)

func main() {
	// CREATE EXTENSION postgis;
	shouldDrop := flag.Bool("shouldDrop", false, "drop tables if existing before migration")
	flag.Parse()
	log.Printf("flag: shouldDrop=%t", *shouldDrop)
	var db *store.DB
	app := fx.New(
		settings.Module,
		fx.Provide(logging.NewZapLogger),
		storeFx.Module,
		fx.Populate(&db),
	)
	err := app.Start(context.Background())
	if err != nil {
		panic(err)
	}
	models := []interface{}{
		&model.User{},
		&model.Project{},
	}
	if *shouldDrop {
		if err := db.Migrator().DropTable(models...); err != nil {
			log.Fatal("Migrator drop table failed")
		}
	}
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatal("Auto migrate failed: ", err)
	}
}
