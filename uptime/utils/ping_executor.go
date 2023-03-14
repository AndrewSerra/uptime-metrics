package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

type PingDuration int64

type PingConfig struct {
	Interval time.Duration
}

type PingExecutor struct {
	Addr    string
	Channel chan bool
}

func (p PingExecutor) Start(ctx context.Context, interval time.Duration, wg *sync.WaitGroup) {
	host, ok := ctx.Value("DBHost").(string)

	if !ok {
		fmt.Println(errors.New("Database host is not given"))
		os.Exit(1)
	}

	token, ok := ctx.Value("DBAccessToken").(string)

	if !ok {
		fmt.Println(errors.New("Database token is not given"))
		os.Exit(2)
	}

	influxClient := influxdb.NewClient(host, token)
	defer influxClient.Close()

	ticker := time.NewTicker(interval)

	for {
		select {
		case keepRunning := <-p.Channel:
			if keepRunning == false {
				fmt.Printf("Ending ping operation for %+v\n", p.Addr)
				ticker.Stop()
				wg.Done()
			}

		case <-ticker.C:
			dur, err := Ping(p.Addr)

			if err != nil {
				log.Fatal(err)
				os.Exit(4)
			}

			err = p.saveDuration(ctx, &influxClient, dur)

			if err != nil {
				fmt.Println(err)
				os.Exit(5)
			}
		}
	}
}

func (p PingExecutor) End(wg *sync.WaitGroup) {
	fmt.Printf("Ending ping operation for %+v\n", p.Addr)
	p.Channel <- true
	wg.Done()
	return
}

func (p PingExecutor) saveDuration(ctx context.Context, client *influxdb.Client, duration *time.Duration) error {
	org, ok := ctx.Value("DBOrg").(string)

	if !ok {
		return errors.New("Database organization is not given")
	}

	bucketName, ok := ctx.Value("DBBucketName").(string)

	if !ok {
		return errors.New("Database bucket name is not given")
	}

	influxWriteAPI := (*client).WriteAPIBlocking(org, bucketName)

	point := influxdb.NewPoint(
		"ping",
		map[string]string{
			"unit":   "nanosecond",
			"domain": p.Addr,
		},
		map[string]interface{}{
			"duration": duration.Nanoseconds(),
		},
		time.Now(),
	)

	err := influxWriteAPI.WritePoint(context.Background(), point)

	if err != nil {
		return err
	}

	return nil
}
