package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"time"
)

type Config struct {
	Address                 string        `short:"a" long:"address" env:"RUN_ADDRESS" default:"localhost:8090" description:"Server host address"`
	LogLevel                string        `short:"l" long:"log" env:"LOG_LEVEL" default:"INFO" description:"Log Level"`
	DatabaseConnection      string        `short:"d" long:"database" env:"DATABASE_URI" description:"Database connection string"`
	JWTSecret               string        `short:"j" long:"jwt" env:"JWT_SECRET" default:"rabbit_Hole" description:"Database connection string"`
	ReportIntervalInSeconds int           `short:"i" long:"interval" env:"REPORT_INTERVAL" default:"10" description:"Frequency (in seconds) for sending requests to the accrual server"`
	ReportInterval          time.Duration `long:"-" description:"Derived duration from ReportIntervalInSeconds"`
	AccrualSystemAddress    string        `short:"r" long:"accrual" env:"ACCRUAL_SYSTEM_ADDRESS" default:"http://localhost:8080" description:"Accrual system host address"`
}

func NewConfig(cliArgs []string) (*Config, error) {
	config := &Config{}
	parser := flags.NewParser(config, flags.Default)

	_, err := parser.ParseArgs(cliArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config.ReportInterval = time.Duration(config.ReportIntervalInSeconds) * time.Second

	return config, nil
}
