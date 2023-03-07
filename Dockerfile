# Build stage
FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN apk update && \
    apk add --no-cache git && \
    apk add --no-cache build-base && \
    go mod download && \
    CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o forum .

# Final stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/forum .
COPY --from=builder /app/data/database.db ./data/database.db
COPY --from=builder /app/utils/cert /app/utils/cert
COPY --from=builder /app/templates /app/templates

EXPOSE 443
CMD ["./forum"]
