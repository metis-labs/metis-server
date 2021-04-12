# Stage 1: build binary
# Start from the latest golang base image
FROM golang:1-buster AS builder

# Add Maintainer Info
LABEL maintainer="MaaS Developers <dl_maas_dev@navercorp.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the metis
RUN make build

# Stage 2: copy binary
FROM debian:buster-slim

# Get and place binary to /bin
COPY --from=builder /app/bin/metis /bin/

# Expose port 10118 to the outside world
EXPOSE 10118

# Define default entrypoint.
ENTRYPOINT ["metis"]
