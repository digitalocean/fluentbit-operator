ARG FLUENTBIT_VERSION=2.2.2

FROM golang:1.21 as builder

WORKDIR /go/src/app

# Download modules first to allow for efficient cache usage.
COPY go.mod go.sum .
RUN go mod download

COPY cmd ./cmd
COPY api ./api
RUN CGO_ENABLED=0 go build -o /tmp/fluent-bit-watcher ./cmd/fluent-bit-watcher/...

FROM fluent/fluent-bit:${FLUENTBIT_VERSION}

COPY conf/fluent-bit.conf conf/parsers.conf /fluent-bit/etc/
COPY --from=builder /tmp/fluent-bit-watcher /fluent-bit/bin/fluent-bit-watcher

# Force creation of the /fluent-bit/config dir. `mkdir` doesn't work in this base image.
COPY conf/empty.conf /fluent-bit/config/fluent-bit.conf

# Entry point
ENTRYPOINT ["/fluent-bit/bin/fluent-bit-watcher"]
