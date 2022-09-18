package config

type Config struct {
	YouTubeAPIKey       string `env:"YOUTUBE_API_KEY"`
	YouTubeVideoQuery   string `env:"YOUTUBE_VIDEO_QUERY"`
	YouTubePollInterval int    `env:"YOUTUBE_POLL_INTERVAL"`

	PostgresHost string `env:"PGHOST,default=localhost"`
	PostgresPort string `env:"PGPORT,default=5432"`
	PostgresUser string `env:"PGUSER"`
	PostgresPass string `env:"PGPASSWORD"`
	PostgresDB   string `env:"PGDB"`

	ServerPort string `env:"PORT,default=8080"`
	ServerHost string `env:"HOST,default=localhost"`
}
