package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type Config struct {
	Address            string `short:"a" long:"address" env:"RUN_ADDRESS" default:"localhost:8090" description:"Server host address"`
	LogLevel           string `short:"l" long:"log" env:"LOG_LEVEL" default:"INFO" description:"Log Level"`
	DatabaseConnection string `short:"d" long:"database" env:"DATABASE_URI" description:"Database connection string"`
	JWTSecret          string `short:"j" long:"jwt" env:"JWT_SECRET" default:"rabbit_Hole" description:"Database connection string"`
}

func NewConfig(cliArgs []string) (*Config, error) {
	config := &Config{}
	parser := flags.NewParser(config, flags.Default)

	_, err := parser.ParseArgs(cliArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return config, nil
}
