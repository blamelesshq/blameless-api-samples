package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var configOncer sync.Once
var config *Config

type Prometheus struct {
	Host string
	Port int
}

type Ingest struct {
	Backfill int
	Period   int
	Step     int
}

type Blameless struct {
	Host      string
	Port      int
	AuthToken string
	OrgId     int
}

type Http struct {
	RequestTimeout int
}

type Config struct {
	Prometheus Prometheus
	Ingest     Ingest
	Blameless  Blameless
	Http       Http
}

type ConfigClient interface {
	Environment() *Config
}

func Environment() *Config {
	configOncer.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Unable to read in config")
		}

		config = &Config{
			Prometheus: Prometheus{
				Host: viper.GetString("prometheus.host"),
				Port: viper.GetInt("prometheus.port"),
			},
			Ingest: Ingest{
				Backfill: viper.GetInt("ingest.backfill"),
				Period:   viper.GetInt("ingest.period"),
				Step:     viper.GetInt("ingest.step"),
			},
			Blameless: Blameless{
				Host:      viper.GetString("blameless.host"),
				Port:      viper.GetInt("blameless.port"),
				AuthToken: viper.GetString("blameless.authToken"),
				OrgId:     viper.GetInt("blameless.orgId"),
			},
			Http: Http{
				RequestTimeout: viper.GetInt("http.requestTimeout"),
			},
		}
	})

	return config
}
