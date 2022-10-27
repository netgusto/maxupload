package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/gammazero/workerpool"
	"go.uber.org/ratelimit"
)

const WORKER_POOL_SIZE = 200     // number of parallel uploads at most
const CALLS_PER_5_MINUTES = 9999 // default is 1200

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cf_api_token := mustGetEnv("CLOUDFLARE_API_TOKEN")
	cf_account_id := mustGetEnv("CLOUDFLARE_ACCOUNT_ID")

	workerRl := ratelimit.New(CALLS_PER_5_MINUTES-100, ratelimit.Per(time.Minute*5)) // - 100 for safety
	wp := workerpool.New(WORKER_POOL_SIZE)

	api, err := cloudflare.NewWithAPIToken(cf_api_token, cloudflare.UsingRateLimit(CALLS_PER_5_MINUTES/(5*60)))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	start := time.Now()
	totalBytes := int64(0)

	nbRunningWorkers := int64(0)

	// import 5000 1.4MB png
	for i := 0; i < 5000; i++ {

		workerRl.Take() // This is a blocking call. Honors our worker dispatch the rate limit

		i := i

		wp.Submit(func() {
			atomic.AddInt64(&nbRunningWorkers, 1)
			defer atomic.AddInt64(&nbRunningWorkers, -1)

			file, err := os.Open("goodboy.png")
			if err != nil {
				panic(err)
			}

			fstat, err := file.Stat()
			if err != nil {
				panic(err)
			}

			uploadStart := time.Now()

			img, err := api.UploadImage(ctx, cf_account_id, cloudflare.ImageUploadRequest{
				File: file,
				Name: fmt.Sprintf("Good boy #%v", i),
			})

			if err != nil {
				panic(err)
			}

			totalBytes += fstat.Size()

			secondsSinceStart := time.Since(start).Seconds()
			fmt.Printf(
				"Upload #%v: %v; took %vms; rate: %.1f/s; throughput: %.1f MB/s; total sent: %.1f MB; %v workers\n",
				i,
				img.ID,
				time.Since(uploadStart).Milliseconds(),
				float64(i)/secondsSinceStart,
				float64(totalBytes/1024/1024)/secondsSinceStart,
				float64(totalBytes)/1024/1024,
				math.Max(0, float64(nbRunningWorkers)),
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
