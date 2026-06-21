FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /bin/output ./cmd

FROM scratch

# 1. Cloud Run relies heavily on TLS for external APIs. Copy certs from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 2. Copy timezone data so your CRUD app handles timezones correctly
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 3. Create a secure, non-root user footprint
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary directly to the root
COPY --from=builder /bin/output /output

EXPOSE 8080

ENTRYPOINT ["/output"]
