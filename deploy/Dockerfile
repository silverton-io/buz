FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY pkg pkg
COPY .VERSION .VERSION
RUN go build -ldflags="-X main.VERSION=$(cat .VERSION)" -o honeypot pkg/*.go
RUN rm -rf pkg

CMD [ "./honeypot" ]
