package config

type Config struct {
	YouTubeAPIKey string `env:"YOUTUBE_API_KEY"`

	PostgresHost string `env:"PGHOST,default=localhost"`
	PostgresPort string `env:"PGPORT,default=5432"`
	PostgresUser string `env:"PGUSER"`
	PostgresPass string `env:"PGPASSWORD"`
	PostgresDB   string `env:"PGDB"`
}
