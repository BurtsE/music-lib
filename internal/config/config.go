package config

type Config struct {
	Postgres Postgres
	Port     string
}

type Postgres struct {
	Username string
	Password string
	Database string
	Hostname string
	Port     int
}
