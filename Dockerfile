FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY pkg pkg
COPY config.yml config.yml
RUN go build -o gosnowplow pkg/*.go
RUN rm -rf pkg

CMD [ "./gosnowplow" ]
