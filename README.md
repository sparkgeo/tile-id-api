# Tile ID API

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

By default the API binds to port 8080. This can be changed by setting an environment variable `TILE_ID_LISTEN_PORT` to a suitable port:
```
go build && TILE_ID_LISTEN_PORT=8129 ./tile-id-api
```
