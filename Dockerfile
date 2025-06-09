# --- Base image with Go
FROM golang:1.23-bullseye AS builder

# Set environment for Go modules
ENV CGO_ENABLED=1
ENV GO111MODULE=on

# Create working directory
WORKDIR /opt/krasiot-sensor

# Copy Go modules first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go binary
RUN go build -o main .

# --- Final image for runtime
FROM debian:bullseye-slim

# Install Oracle Instant Client dependencies
RUN apt-get update \
  && apt-get install -y --no-install-recommends libaio1 ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# Copy Instant Client before setting up linker
COPY ./instantclient_19_10 /opt/oracle/instantclient_19_10

# Symlink and configure linker after copying libraries
RUN ln -s /opt/oracle/instantclient_19_10 /opt/oracle/instantclient \
  && echo /opt/oracle/instantclient > /etc/ld.so.conf.d/oracle-instantclient.conf \
  && ldconfig

# Environment variables for Oracle
ENV LD_LIBRARY_PATH=/opt/oracle/instantclient
ENV TNS_ADMIN=/opt/krasiot-sensor/wallet

# Create working directory
WORKDIR /opt/krasiot-sensor

# Copy the built binary
COPY --from=builder /opt/krasiot-sensor/main .

RUN chmod +x main

# Entrypoint
CMD ["./main"]
