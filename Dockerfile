# Use official Golang image
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Install git (needed for go mod sometimes)
RUN apk update && apk add --no-cache git

# Copy go mod and sum
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the source
COPY . .

# Build the Go app
RUN go build -o main .

# Run it
CMD ["/app/main"]
