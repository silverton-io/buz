FROM golang:1.17-alpine

WORKDIR /app
ENV GOOGLE_APPLICATION_CREDENTIALS=/config/dp-dev.json

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY pkg pkg
COPY config.yml config.yml
RUN go build -o gosnowplow pkg/*.go

COPY dp-dev.json /config/dp-dev.json

CMD [ "./gosnowplow" ]