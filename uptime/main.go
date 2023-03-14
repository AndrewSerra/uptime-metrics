package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sm/uptime/utils"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type InfluxDatabaseConfig struct {
	Bucket string `yaml:"bucket"`
	Host   string `yaml:"host"`
	Token  string `yaml:"token"`
	Org    string `yaml:"org"`
}

type PingConfig struct {
	Amount time.Duration `yaml:"amount"`
}

type UptimeConfig struct {
	Database InfluxDatabaseConfig `yaml:"database"`
	Time     PingConfig           `yaml:"time"`
}

type PingData struct {
	duration int64
}

func parseYaml(data []byte) (*UptimeConfig, error) {
	uptimeConfig := &UptimeConfig{}

	err := yaml.Unmarshal(data, uptimeConfig)

	if err != nil {
		return nil, err
	}

	return uptimeConfig, nil
}

func createContext(configs *UptimeConfig) context.Context {
	ctx := context.Background()

	// Add Database Info to context
	ctx = context.WithValue(ctx, "DBBucketName", configs.Database.Bucket)
	ctx = context.WithValue(ctx, "DBHost", configs.Database.Host)
	ctx = context.WithValue(ctx, "DBAccessToken", configs.Database.Token)
	ctx = context.WithValue(ctx, "DBOrg", configs.Database.Org)

	return ctx
}

func main() {
	// Check domain count passed to service
	if urls := os.Args[1:]; len(urls) < 1 {
		log.Fatal("There has to be at least one address to ping.")
		os.Exit(1)
	}

	fmt.Println("Parsing config file for uptime service...")

	yamlFile, err := os.ReadFile("./config/uptime.yaml")

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	configs, err := parseYaml(yamlFile)

	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}

	ctx := createContext(configs)
	urls := os.Args[1:]
	var wg sync.WaitGroup

	chanCoordinator := utils.ChannelCoordinator{}

	for _, url := range urls {
		uc, err := chanCoordinator.Add(url)

		if err != nil {
			log.Fatal(err)
			os.Exit(4)
		}

		pinger := &utils.PingExecutor{
			Addr:    url,
			Channel: uc.C,
		}

		wg.Add(1)
		go pinger.Start(ctx, configs.Time.Amount, &wg)
	}

	wg.Wait()
}
