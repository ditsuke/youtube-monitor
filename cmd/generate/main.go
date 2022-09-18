package main

import (
	"context"
	"github.com/ditsuke/youtube-focus/config"
	"github.com/ditsuke/youtube-focus/internal/yt"
	"github.com/ditsuke/youtube-focus/store"
	"github.com/sethvargo/go-envconfig"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
)

// prepareDb prepares a postgres store with the req table
func prepareDb(db *gorm.DB) error {
	return db.AutoMigrate(yt.VideoFull{})
}

func main() {
	cfg := config.Config{}
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatalln(err)
	}

	db := store.GetDB(store.GetDSNFromConfig(cfg))

	// prepare table
	if err := prepareDb(db); err != nil {
		log.Fatalln(err)
	}

	// generate model(s)
	g := gen.NewGenerator(gen.Config{OutPath: "model"})
	g.UseDB(db)
	g.GenerateAllTable()
	g.Execute()
}
