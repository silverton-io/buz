FROM golang:1.17 AS build

WORKDIR /app

COPY . /app/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.VERSION=$(cat .VERSION)" -o honeypot ./cmd/*.go

FROM busybox

LABEL maintainer="Jake Thomas <jake@bostata.com>"
LABEL org.opencontainers.image.description "A lightweight, Snowplow-compatible streaming event collection system."

WORKDIR /app
COPY --from=build /app/.VERSION .
COPY --from=build /app/honeypot .
COPY --from=build /etc/ssl/certs /etc/ssl/certs

EXPOSE 8080

ENTRYPOINT [ "./honeypot" ]
