package main

import (
	"fmt"
	"log"
	"os"
	"sm/uptime/utils"

	yaml "gopkg.in/yaml.v3"
)

type InfluxDatabaseConfig struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
	Org      string `yaml:"org"`
}

type InfluxConfig struct {
	Database InfluxDatabaseConfig `yaml:"database"`
}

func parseYaml(data []byte) (*InfluxDatabaseConfig, error) {
	influxConfig := &InfluxConfig{}

	err := yaml.Unmarshal(data, &influxConfig)

	if err != nil {
		return nil, err
	}

	return &influxConfig.Database, nil
}

func main() {
	fmt.Println("Parsing config file for influxdb...")

	yamlFile, err := os.ReadFile("./config/influxdb.yaml")

	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	configs, err := parseYaml(yamlFile)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", configs)

	dur, err := utils.Ping("google.com")

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	fmt.Printf("Duration: %+d\n", dur)
}
