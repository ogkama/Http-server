package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
  Address           string       
}

func MustLoad(cfgPath string, cfg any) {
  if cfgPath == "" {
    log.Fatal("Config path is not set")
  }

  if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
    log.Fatalf("config file does not exist by this path: %s", cfgPath)
  }

  if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
    log.Fatalf("error reading config: %s", err)
  }
}

type AppFlags struct {
  ConfigPath        string
}

func ParseFlags() AppFlags {
  configPath := flag.String("config", "", "Path to config")
  flag.Parse()
  return AppFlags{
    ConfigPath:        *configPath,
  }
}

type RabbitMQ struct {
	Host 	 	string `yaml:"host"`
	Port 	 	uint16 `yaml:"port"`
	QueueName 	string `yaml:"queue_name"`
}

type AppConfig struct {
	RabbitMQ   `yaml:"rabbit_mq"`
	HTTPConfig `yaml:"http"`
}