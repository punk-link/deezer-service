# deezer-service
Deezer API implementation for release syncing


## Environment Variables

```json
{
  "HostSettings": {
  	"Address": "127.0.0.1",
  	"Port": "8080",
  	"Mode": "release"
  },
  "JaegerSettings": {
  	"Endpoint": ""
  },
  "NatsSettings": {
  	"Endpoint": ""
  },
  "SentryDsn": ""
}
```

## Metrics

|Name|Type|
|-|-|
|release_urls_in_process|gauge|
|urls_received_total|counter|
|urls_processed|counter|