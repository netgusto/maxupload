# maxupload

Demonstration of upload to Cloudflare Images, maximizing upload rate while remaining below rate limit (default is 1200 api calls per 5 minutes, ie an avg of 4/s).

Limits set to 200 workers / 9999 calls per 5 minutes (change this in the code to adapt values to your case).

Bundled image is https://unsplash.com/photos/t8T1xTFqAXY by Michael G.

## usage

```bash
CLOUDFLARE_API_TOKEN="the-token-here" CLOUDFLARE_ACCOUNT_ID="images-account-id-here" go run .
```

Presents the uploaded image ID, duration for image, overall upload rate/s since beginning, and overall throughput/s since beginning.

Tested from a Digital Ocean droplet in Europe (Frankfurt), uploads reached 24/s limited on 100 workers / 9999 calls per 5 minute (theoretical limit is 33/s with this rate limiting).

```
Upload #359: 84d1dcca-86fe-4c41-20dc-ce3a6c1b9300; took 3536ms; rate: 23.3/s; throughput: 40.3 MB/s; 100 workers
Upload #371: c744dcbb-99a7-43d8-fe29-019d31a36400; took 3313ms; rate: 24.1/s; throughput: 40.4 MB/s; 100 workers
Upload #373: 95ce1b03-1b78-4c93-1011-b639de9b6a00; took 3303ms; rate: 24.1/s; throughput: 40.3 MB/s; 100 workers
Upload #372: 543aa6f7-f7d3-4828-4bac-2fea53fb5b00; took 3320ms; rate: 24.0/s; throughput: 40.4 MB/s; 100 workers
Upload #370: 9fe63213-9743-414f-31fb-79648ba11100; took 3391ms; rate: 23.9/s; throughput: 40.5 MB/s; 100 workers
Upload #366: ca64f7e5-a1ff-4714-c44c-7edd836c4f00; took 3569ms; rate: 23.6/s; throughput: 40.5 MB/s; 100 workers
```

## 🚨 Note on rate limiting 🚨

It is important to not push uploads to the actual rate limit of the account, as all calls made on this account by other systems are accounted for by Cloudflare, not only those of uploads; pushing it to the max would risk your account being (temporarily) blocked for api hammering.
