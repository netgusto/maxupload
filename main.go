package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/gammazero/workerpool"
	"go.uber.org/ratelimit"
)

const WORKER_POOL_SIZE = 30 // number of parallel uploads at most

func main() {
	cf_api_token := mustGetEnv("CLOUDFLARE_API_TOKEN")
	cf_account_id := mustGetEnv("CLOUDFLARE_ACCOUNT_ID")

	rl := ratelimit.New(1200-50, ratelimit.Per(time.Minute*5)) // 1200 requests every 5 minutes - 50 for safety

	wp := workerpool.New(WORKER_POOL_SIZE)

	api, err := cloudflare.NewWithAPIToken(cf_api_token)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	start := time.Now()
	totalBytes := int64(0)

	// import 5000 images
	for i := 0; i < 5000; i++ {

		rl.Take() // This is a blocking call. Honors the rate limit

		i := i

		wp.Submit(func() {
			file, err := os.Open("goodboy.png")
			if err != nil {
				panic(err)
			}

			fstat, err := file.Stat()
			if err != nil {
				panic(err)
			}

			totalBytes += fstat.Size()
			uploadStart := time.Now()

			time.Sleep(time.Millisecond * 1000)

			img, err := api.UploadImage(ctx, cf_account_id, cloudflare.ImageUploadRequest{
				File: file,
				Name: fmt.Sprintf("Good boy #%v", i),
			})

			if err != nil {
				panic(err)
			}

			secondsSinceStart := time.Since(start).Seconds()
			fmt.Printf(
				"Upload #%v: %v; took %vms; rate: %v/s; throughput: %v MB/s\n",
				i,
				img.ID,
				time.Since(uploadStart).Milliseconds(),
				math.Round(float64(i)/secondsSinceStart),
				math.Round(float64(totalBytes/1024/1024)/secondsSinceStart),
			)
		})
	}
}

func mustGetEnv(name string) string {
	v, isSet := os.LookupEnv(name)
	if !isSet {
		panic(fmt.Sprintf("Missing env var: %v", name))
	}

	return v
}
