# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /src

# CA bundle is copied into the scratch image for libSQL TLS.
RUN apk add --no-cache ca-certificates

# Cache module layer separately from source.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static binary: no libc needed, so legacy scratch works as runtime.
RUN CGO_ENABLED=0 GOOS=linux go build \
      -trimpath \
      -ldflags="-s -w" \
      -o /out/meowcounter \
      ./cmd/api

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/meowcounter /meowcounter

# nobody:nogroup
USER 65534:65534

ENV PORT=3000
EXPOSE 3000

ENTRYPOINT ["/meowcounter"]
