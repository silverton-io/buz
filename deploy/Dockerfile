FROM golang:1.19-bullseye AS build

WORKDIR /app

COPY . /app/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o healthcheck ./cmd/healthcheck/*.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.VERSION=$(cat .VERSION)" -o honeypot ./cmd/honeypot/*.go

FROM busybox

LABEL maintainer="Jake Thomas <jake@silverton.io>"
LABEL org.opencontainers.image.description "A lightweight, Snowplow-compatible streaming event collection system."

WORKDIR /app
COPY --from=build /app/.VERSION .
COPY --from=build /app/honeypot .
COPY --from=build /app/healthcheck .
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /app/examples/quickstart/honeypot/quickstart.conf.yml /etc/honeypot/config.yml

EXPOSE 8080
ENV HONEYPOT_CONFIG_PATH=/etc/honeypot/config.yml

ENTRYPOINT [ "./honeypot" ]
