# SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

# Stage 1: Build the Go app
FROM golang:1.24 AS build

WORKDIR /go/src
# Copy the source code into the container
COPY . .
RUN go mod download

# Build the application
RUN make watcher

# Stage 2: Produce the runtime image
FROM quay.io/brancz/kube-rbac-proxy:v0.18.2 AS watcher

# Copy the binary from the build stage
COPY --from=build /go/src/build/watcher /usr/bin/watcher

# Command to run the application
ENTRYPOINT ["/usr/bin/watcher"]
