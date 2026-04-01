FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(git describe --tags --always 2>/dev/null || echo dev)" -o /bin/roster ./cmd/roster/
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl
COPY --from=builder /bin/roster /usr/local/bin/roster
ENV PORT="8930" DATA_DIR="/data" ROSTER_LICENSE_KEY=""
EXPOSE 8930
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD curl -sf http://localhost:8930/health || exit 1
ENTRYPOINT ["roster"]
