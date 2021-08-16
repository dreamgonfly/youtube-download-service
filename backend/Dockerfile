# Use the offical golang image to create a binary.
FROM golang:1.16-buster as builder

ENV LANG C.UTF-8

# Create and change to the app directory.
WORKDIR /app

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
    curl

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o serve

# Use the official Debian slim image for a lean production container.
FROM debian:buster-slim

ENV LANG C.UTF-8

WORKDIR /app

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    curl \
    python3 \
    aria2 \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*
RUN ln -s /usr/bin/python3 /usr/bin/python

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/serve .
COPY --from=builder /app/.credentials /credentials
COPY --from=builder /app/configs configs

# Install youtube-dl
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl

# Run the web service on container startup.
CMD ["./serve"]
