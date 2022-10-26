# maxupload

Demonstration of upload to Cloudflare Images, maximizing upload rate while remaining below rate limit (default is 1200 api calls per 5 minutes, ie an avg of 4/s).

Bundled image is https://unsplash.com/photos/t8T1xTFqAXY by Michael G.

## usage

```bash
CLOUDFLARE_API_TOKEN="the-token-here" CLOUDFLARE_ACCOUNT_ID="images-account-id-here" go run .
```

Presents the uploaded image ID, duration for image, overall upload rate/s since beginning, and overall throughput/s since beginning.
