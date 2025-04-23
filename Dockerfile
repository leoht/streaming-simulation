FROM golang:1.24


WORKDIR /app

COPY . .

RUN go mod download 
# RUN go get "github.com/confluentinc/confluent-kafka-go/v2/kafka"

# ENV


RUN CGO_ENABLED=1 go build -o /app/cmd/web/main -ldflags '-linkmode external -w -extldflags "-static"' /app/cmd/web/main.go

EXPOSE 8080

# Set environment variable for Gin mode
# ENV GIN_MODE=release

# Run the executable
CMD ["/app/cmd/web/main"]