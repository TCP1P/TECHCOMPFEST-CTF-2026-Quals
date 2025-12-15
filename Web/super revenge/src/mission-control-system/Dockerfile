FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG GO_BUILD_TARGET=.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /tmp/server "${GO_BUILD_TARGET}"

FROM alpine:3.20 AS runner
WORKDIR /app
RUN apk add --no-cache ca-certificates && adduser -D -g '' appuser
COPY --from=builder /tmp/server /app/server
COPY .env .env
USER appuser

ENV PORT=3002
ENV ENV=production

EXPOSE 3002
CMD ["./server"]
