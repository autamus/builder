# Start from the latest golang base image
FROM ghcr.io/autamus/go:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o builder .

# Start again with minimal envoirnment.
FROM ghcr.io/spack/ubuntu-bionic:latest

RUN apt-get update && \
    apt-get install -y ca-certificates gringo

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder /app/builder /app/builder
COPY spack/images.json /opt/spack/lib/spack/spack/container/images.json

ENV PATH=/opt/spack/bin:$PATH

# Command to run the executable
ENTRYPOINT ["/app/builder"]
