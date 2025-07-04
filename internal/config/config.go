package config

import (
	"github.com/jessevdk/go-flags"
	"time"
)

type Config struct {
	Address                   string        `short:"a" long:"address" env:"RUN_ADDRESS" default:"localhost:8090" description:"Server host address"`
	LogLevel                  string        `short:"l" long:"log" env:"LOG_LEVEL" default:"INFO" description:"Log Level"`
	DatabaseConnection        string        `short:"d" long:"database" env:"DATABASE_URI" description:"Database connection string"`
	JWTSecret                 string        `short:"j" long:"jwt" env:"JWT_SECRET" default:"rabbit_Hole" description:"Database connection string"`
	ReportIntervalInSeconds   int           `short:"i" long:"interval" env:"REPORT_INTERVAL" default:"10" description:"Frequency (in seconds) for sending requests to the accrual server"`
	ReportInterval            time.Duration `long:"-" description:"Derived duration from ReportIntervalInSeconds"`
	AccrualSystemAddress      string        `short:"r" long:"accrual" env:"ACCRUAL_SYSTEM_ADDRESS" default:"http://localhost:8080" description:"Accrual system host address"`
	GracefulShutdownInSeconds int           `short:"s" long:"shutdown" env:"SHUTDOWN_INTERVAL" default:"30" description:"Frequency (in seconds) for graceful shutdown"`
	GracefulShutdownInterval  time.Duration `description:"Derived duration from GracefulShutdownInSeconds"`
}

func NewConfig(cliArgs []string) (*Config, error) {
	config := &Config{}
	parser := flags.NewParser(config, flags.Default)

	_, err := parser.ParseArgs(cliArgs)
	if err != nil {
		panic(err)
	}

	config.ReportInterval = time.Duration(config.ReportIntervalInSeconds) * time.Second
	config.GracefulShutdownInterval = time.Duration(config.GracefulShutdownInSeconds) * time.Second
	return config, nil
}
