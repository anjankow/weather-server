# Weather forecast server
A simple HTTP server aggregating weather forecast from multiple forecast providers.

The server and all its dependencies are initialized in `cmd/api-server/main.go`.

## Run
1. Navigate to the project root.
2. Copy .env.example to .env
3. Set `WEATHER_API_KEY`
4. Run
```sh
go run cmd/api-server/main.go
```

## Test
```sh
curl 'localhost:8099/weather?longitude=44&latitude=11'
```

# API
## GET `/weather`
Required query params:
- `longitude` range `-180` to `180`
- `latitude` range `-90` to `90`

## GET `/-/healthy`

