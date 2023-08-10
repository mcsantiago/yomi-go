# Start from the official Golang base image
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

COPY ./data/JMdict_e /data/JMdict_e

# Download all the dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/tokenizer-apiserver/apiserver.go

# Start a new stage from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the data files from the previous stage
COPY --from=builder /data/JMdict_e /data/JMdict_e

# Expose port to the outside world (if applicable)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

