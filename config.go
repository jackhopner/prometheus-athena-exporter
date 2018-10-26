package main

import (
	"flag"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	ListenAddress string         `yaml:"listen-address"`
	Tenants       []tenantConfig `yaml:"tenants"`
	SSMPrefix     string         `yaml:"ssm-prefix"`
	AWSRegionID   string         `yaml:"aws-region-id"`
	AWSAccountID  string         `yaml:"aws-account-id"`

	Metrics []metric `yaml:"metrics"`
}

type tenantConfig struct {
	Tenant string `yaml:"tenant"`
	DBName string `yaml:"db-name"`
}

type metric struct {
	Name              string        `yaml:"name"`
	Query             string        `yaml:"query"`
	QueryValueColumns []string      `yaml:"query-value-columns"`
	QueryInterval     time.Duration `yaml:"query-interval"`
	ExlcudeDBs        []string      `yaml:"exclude-dbs"`
	IncludeDBS        []string      `yaml:"include-dbs"`
}

func loadConfig() config {
	var configFile string
	result := config{}

	flags := flag.NewFlagSet("fs", flag.PanicOnError)
	flags.StringVar(&configFile, "config", "", "YAML config file")
	flags.Parse(os.Args[1:])

	if configFile == "" {
		flags.Usage()
		log.Panicf("config file is required")
	}

	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.WithError(err).Panicf("cannot read config file")
	}

	if err = yaml.Unmarshal(raw, &result); err != nil {
		log.WithError(err).Panicf("failed to unmarshal contents of config file")
	}

	log.WithField("config", result).Info("loaded config")

	return result
}
