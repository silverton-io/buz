FROM golang:1.17-alpine


LABEL org.opencontainers.image.description DESCRIPTION

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY cmd cmd
COPY pkg pkg
COPY .VERSION .VERSION

RUN go mod download
RUN go build -ldflags="-X main.VERSION=$(cat .VERSION)" -o honeypot ./cmd/*.go
RUN rm -rf pkg

ENTRYPOINT [ "./honeypot" ]
