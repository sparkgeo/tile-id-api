FROM golang:1.21.6-alpine3.19 as builder

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build


FROM alpine:3.19 

COPY --from=builder /src/tile-id-api /
EXPOSE 8080
CMD ["/tile-id-api"]
