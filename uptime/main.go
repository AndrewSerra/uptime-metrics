package main

import (
	"fmt"
	"log"
	"os"
	"sm/uptime/utils"
	"time"

	yaml "gopkg.in/yaml.v3"
)

// type PingIntervalUnit string

// const (
// 	SECOND PingIntervalUnit = "second"
// 	MINUTE PingIntervalUnit = "minute"
// 	HOUR   PingIntervalUnit = "hour"
// 	DAY    PingIntervalUnit = "day"
// )

type InfluxDatabaseConfig struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
	Org      string `yaml:"org"`
}

type PingConfig struct {
	// Unit   PingIntervalUnit `yaml:"unit"`
	Amount time.Duration `yaml:"amount"`
}

type InfluxConfig struct {
	Database InfluxDatabaseConfig `yaml:"database"`
	Time     PingConfig           `yaml:"time"`
}

type PingData struct {
}

func parseYaml(data []byte) (*InfluxConfig, error) {
	influxConfig := &InfluxConfig{}

	err := yaml.Unmarshal(data, influxConfig)

	if err != nil {
		return nil, err
	}

	return influxConfig, nil
}

func main() {
	fmt.Println("Parsing config file for uptime service...")

	yamlFile, err := os.ReadFile("./config/uptime.yaml")

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

	ticker := time.NewTicker(time.Duration(configs.Time.Amount))
	defer ticker.Stop()

	for _ = range ticker.C {
		// fmt.Println("Current time: ", tick)
		dur, err := utils.Ping("google.com")

		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}

		fmt.Printf("Duration: %+d\n", dur)
	}
}
