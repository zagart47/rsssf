package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type cfg struct {
	RSS           []string `json:"rss"`
	requestPeriod int      `json:"request_period"`
	Timeout       time.Duration
	Postgres      struct {
		DSN      string
		host     string `json:"host" env-default:"postgres"`
		port     string `json:"port" env-default:"5432"`
		dbname   string `json:"db_name" env-default:"postgres"`
		username string `json:"user_name" env-default:"postgres"`
		pwd      string `json:"pwd" env-default:"postgres"`
	} `json:"postgres"`
}

func config() cfg {
	configs := cfg{}
	if err := cleanenv.ReadConfig("./config/config.json", &configs); err != nil {
		log.Println("cannot read configs")
	}
	configs.Postgres.DSN = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", configs.Postgres.dbname, configs.Postgres.username,
		configs.Postgres.pwd, configs.Postgres.host, configs.Postgres.port, configs.Postgres.dbname)
	configs.Timeout = time.Duration(configs.requestPeriod) * time.Minute
	return configs
}

var Configs = config()
