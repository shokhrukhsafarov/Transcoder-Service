# Build stage
FROM golang:1.22 as builder

# Set up the workspace
RUN mkdir -p $GOPATH/src/gitlab.udevs.io/ucode/ucode_go_transcoder_service
WORKDIR $GOPATH/src/gitlab.udevs.io/ucode/ucode_go_transcoder_service

# Copy the local package files to the container's workspace
COPY . ./

# Install dependencies and build
RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    go mod vendor && \
    make build && \
    mv ./bin/ucode_go_transcoder_service / && \
    mv ./config /

# Final stage
FROM alpine:latest

# Install postgresql-client to provide pg_isready
RUN apk add --no-cache postgresql-client

# Copy the binary and config from the builder stage
COPY --from=builder /ucode_go_transcoder_service .
COPY --from=builder /config ./config

# Set the entrypoint
ENTRYPOINT ["/ucode_go_transcoder_service"]