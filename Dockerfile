FROM golang:1.18
WORKDIR /app
COPY . .
RUN go build -ldflags "-X main.version=1.0.0" -o consumerapp ./consumer
RUN go build -ldflags "-X main.version=1.0.0" -o producerapp ./producer
CMD ["./consumerapp", "./producerapp"]
