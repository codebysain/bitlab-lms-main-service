FROM golang:1.23.8

WORKDIR /app

# Install dependencies including PostgreSQL client tools
RUN apt-get update && \
    apt-get install -y git curl postgresql-client && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install goose for DB migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy go mod files and download deps
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Entrypoint scripts
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Build binary
WORKDIR /app/cmd
RUN go build -o /app/main .

# Set default command
ENTRYPOINT ["/entrypoint.sh"]
