FROM golang:1.23.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY entrypoint.sh /

RUN CGO_ENABLED=0 GOOS=linux go build -o /sp-build

ENTRYPOINT ["/entrypoint.sh"]