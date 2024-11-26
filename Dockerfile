FROM golang:1.23.3-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /sp-build

CMD ["/sp-build"]