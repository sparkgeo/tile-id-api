# Tile ID API

[![Validate and Test](https://github.com/sparkgeo/tile-id-api/actions/workflows/test.yml/badge.svg)](https://github.com/sparkgeo/tile-id-api/actions/workflows/test.yml)

Simple HTTP API to support map tile identification across a number of tiling schemes, written in Go.

Valid requests produce a tile response, with the tile's ID in multiple tiling schemes shown in the image. The same tile IDs are also provided as HTTP response headers alongside the tile's bounds in WGS84 Longitude / Latitude.

## Supported Tiling Schemes

The following tiling schemes are supported. The example request URLs return tiles at the same location.

- xyz (or zxy)
    - Example 1: `http://localhost:8080/xyz/12/234/3245`
    - Example 2: `http://localhost:8080/xyz/0/0/0`
- tms (EPSG:3857 only)
    - Example 1: `http://localhost:8080/tms/12/234/850`
    - Example 2: `http://localhost:8080/tms/0/0/0`
- quadkey
    - Example 1: `http://localhost:8080/quadkey/220031303212`
    - Example 2: `http://localhost:8080/quadkey/`

## Image Encoding and Opacity

By default tiles are encoded as PNGs with 30% opacity. Both PNG and JPEG encoding are available.

Image encoding can be controlled with a `.png`, `.jpg`, or `.jpeg` suffix on a tile request:
```
http://localhost:8080/xyz/12/234/3245.jpg
http://localhost:8080/xyz/12/234/3245.png
```

JPEG encoding does not support opacity and JPEG tiles will always have 100% opacity. Non-permitted encoding suffixes will be rejected.

With PNG encoding opacity can be changed with either a query string parameter or an HTTP request header:
```
http://localhost:8080/xyz/12/234/3245?opacityPercent=80
curl http://localhost:8080/xyz/12/234/3245 -H "X-Opacity-Percent: 80"
```

Only integer opacity values between 0 and 100 are permitted. Non-permitted values will result in default opacity.

## Host Port

By default the API binds to port 8080. This can be changed by setting a flag `server-port` to a suitable port:
```sh
cd tile-id-api/
go build ./cmd/tile-id-api/ && ./tile-id-api --server-port 8123
```

## Log Level

By default the API logs `info` messages. Change this to a supported log level by setting a flag `log-level`. Supported log levels include `panic, fatal, error, warn, info, debug, trace`:
```sh
cd tile-id-api/
go build ./cmd/tile-id-api/ && ./tile-id-api --log-level debug
```

## Docker

The API can be run in a container and by default the container will bind to local port `8123`. This can be changed by setting an environment variable `TILE_ID_LISTEN_PORT` to a suitable port:
```sh
TILE_ID_LISTEN_PORT=9123 scripts/run-in-container.sh
```

## Tests

Unit tests can be executed locally:
```sh
cd tile-id-api/
go test ./... -v
```

Other tests require a container to execute. The following scripts all execute in containers:
```sh
scripts/run-validate-openapi-in-container.sh    # validate the OpenAPI specification
scripts/run-unit-tests-in-container.sh  # execute all unit tests
scripts/run-integration-tests-in-container.sh # execute all integration tests

# this script is executed in CI
scripts/run-tests-in-container.sh   # validate OpenAPI spec, execute unit tests, execute integration tests
```

## MapProxy

A MapProxy configuration is provided to test each of the available tiling schemes as a tile source. This happens as part of integration testing, but if necessary MapProxy can be run in a container. By default MapProxy will bind to local port `8124`. This can be changed by setting an environment variable `MAPPROXY_LISTEN_PORT` to a suitable port: 
```sh
MAPPROXY_LISTEN_PORT=9124 scripts/run-mapproxy-in-container.sh
```

Use `http://localhost:8124/wms?request=GetCapabilities` to access MapProxy's WMS GetCapabilities document. This MapProxy configuration is for testing purposes only, and is not suitable for a production environment.

## Development

Requires Go 1.21.6 for API development.

Requires Python >= 3.8 for integration test development. To configure the local Python environment:
```sh
scripts/dev-init.sh
```