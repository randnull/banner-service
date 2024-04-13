FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o /banner-service ./cmd/banner-app/main.go

CMD ["/banner-service"]