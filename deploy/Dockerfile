FROM golang:1.19-bullseye AS build

WORKDIR /app

COPY . /app/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o healthcheck ./cmd/healthcheck/*.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.VERSION=$(cat .VERSION)" -o buz ./cmd/buz/*.go

FROM busybox

LABEL maintainer="Jake Thomas <jake@silverton.io>"
LABEL org.opencontainers.image.description "A lightweight, Snowplow-compatible streaming event collection system."

WORKDIR /app
COPY --from=build /app/.VERSION .
COPY --from=build /app/buz .
COPY --from=build /app/healthcheck .
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /app/examples/quickstart/buz/quickstart.conf.yml /etc/buz/config.yml

EXPOSE 8080
ENV BUZ_CONFIG_PATH=/etc/buz/config.yml

ENTRYPOINT [ "./buz" ]
