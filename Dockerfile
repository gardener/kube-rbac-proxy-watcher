# Stage 1: Build the Go app
FROM golang:1.23.1 AS build

# Set up the working directory
WORKDIR /src

# Fetch Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
ARG LD_FLAGS
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=0 go build -ldflags="${LD_FLAGS}" -o watcher ./cmd/watcher/main.go

# Stage 2: Produce the runtime image
FROM quay.io/brancz/kube-rbac-proxy:v0.18.1

# Copy the binary from the build stage
COPY --from=build /src/watcher /usr/bin/watcher

# Command to run the application
ENTRYPOINT ["/usr/bin/watcher"]
