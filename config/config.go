package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	YouTubeAPIKeys      []string `env:"YOUTUBE_API_KEYS"`
	YouTubeVideoQuery   string   `env:"YOUTUBE_VIDEO_QUERY,default=game"`
	YouTubePollInterval int      `env:"YOUTUBE_POLL_INTERVAL,default=20"`

	PostgresHost string `env:"PGHOST,default=localhost"`
	PostgresPort string `env:"PGPORT,default=5432"`
	PostgresUser string `env:"PGUSER"`
	PostgresPass string `env:"PGPASSWORD"`
	PostgresDB   string `env:"PGDB"`

	ServerPort string `env:"PORT,default=8080"`
	ServerHost string `env:"HOST,default=localhost"`
}

func (c Config) GetDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.GetDSN()))
}

func (c Config) GetDSN() string {
	return fmt.Sprintf(
		"user=%s password=%s port=%s host=%s dbname=%s",
		c.PostgresUser, c.PostgresPass,
		c.PostgresPort, c.PostgresHost,
		c.PostgresDB,
	)
}
