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

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl \
    && addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /bin/output /app/output

RUN chown app:app /app/output

USER app

EXPOSE 8080

ENTRYPOINT ["/app/output"]
