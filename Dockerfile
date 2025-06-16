FROM golang:1.23.8

WORKDIR /app

# Install dependencies
RUN apt-get update && apt-get install -y git curl

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy everything
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Build binary
WORKDIR /app/cmd
RUN go build -o /app/main .

ENTRYPOINT ["/entrypoint.sh"]


