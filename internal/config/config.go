package config

type Config struct {
	HelperApi string    `env:"HELPER_API"`
	Postgres  *Postgres `json:"postgres"`
	Service   *Service  `json:"service"`
	Host      *Host     `json:"host"`
}
type Host struct {
	Address string `env:"HOST_ADDRESS"`
	Port    int    `env:"HOST_PORT"`
}
type Postgres struct {
	Username string `env:"PG_USER"`
	Password string `env:"PG_PASS"`
	Database string `env:"PG_DB"`
	Hostname string `env:"PG_HOST"`
	Port     int    `env:"PG_PORT"`
	MaxConns int    `json:"max_conns"`
	Sslmode  string `json:"sslmode"`
}
type Service struct {
	PageLimit int `env:"SERVICE_PAGE_LIMIT"`
}
