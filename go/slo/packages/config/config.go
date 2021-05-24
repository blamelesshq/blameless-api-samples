package config

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Prometheus struct {
	Host string
	Port int
}

type Ingest struct {
	Period int
}

type Blameless struct {
	Host string
	Port int
}

type Config struct {
	Prometheus *Prometheus
	Ingest     *Ingest
	Blameless  *Blameless
	Log        *zap.Logger
}

type ConfigClient interface {
	NewConfig() *Config
}

func NewConfig() *Config {
	l := logInit()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		l.Fatal("Unable to read in config")
	}

	p := &Prometheus{
		Host: viper.GetString("prometheus.host"),
		Port: viper.GetInt("prometheus.port"),
	}

	i := &Ingest{
		Period: viper.GetInt("ingest.period"),
	}

	b := &Blameless{
		Host: viper.GetString("blameless.host"),
		Port: viper.GetInt("blameless.port"),
	}

	return &Config{
		Prometheus: p,
		Ingest:     i,
		Blameless:  b,
		Log:        l,
	}
}

func logInit() *zap.Logger {
	errFilter := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	stdFilter := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl > zapcore.DebugLevel
	})

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), stdFilter),
		zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), errFilter),
	)

	return zap.New(core, zap.Fields(zap.String("application", "Blameless SLO")))
}
