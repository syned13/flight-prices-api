FROM golang:1.23-alpine as builder
ENV APP_HOME /go/src/github.com/syned13/flight-prices-api

RUN mkdir -p $APP_HOME
ADD . $APP_HOME
WORKDIR $APP_HOME
RUN mkdir -p build

RUN --mount=type=cache,target=/go/pkg/mod go mod download

RUN go build -o build/main ./cmd/main.go

RUN chmod +x ./build/main

FROM alpine:3.11.3
COPY --from=builder /go/src/github.com/syned13/flight-prices-api .

CMD ["./build/main"]