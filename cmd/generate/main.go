package main

import (
	"context"
	"flag"
	"github.com/ditsuke/youtube-focus/config"
	"github.com/ditsuke/youtube-focus/internal/yt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sethvargo/go-envconfig"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
)

const TSVIndexQuery = `CREATE INDEX IF NOT EXISTS ts_idx ON videos USING GIN (tsv)`

// prepareDb prepares a postgres store with the req table
func prepareDb(db *gorm.DB) error {
	db.Exec(TSVIndexQuery)
	return db.AutoMigrate(yt.VideoFull{})
}

func main() {
	noGen := flag.Bool("no-gen", false, "if true, the db is only prepared")
	flag.Parse()

	cfg := config.Config{}
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatalln(err)
	}

	db, err := cfg.GetDB()
	if err != nil {
		panic(err)
	}

	// prepare table
	if err := prepareDb(db); err != nil {
		log.Fatalln(err)
	}

	log.Println("db prepared")

	if *noGen {
		log.Println("not generating code because of the -no-gen flag")
		return
	}

	// generate model(s)
	g := gen.NewGenerator(gen.Config{OutPath: "model", FieldCoverable: true})
	g.UseDB(db)
	g.GenerateAllTable()
	g.Execute()
	log.Println("models generated")
}
