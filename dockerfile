# Build stage
FROM golang:1.22.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o xyz-multifinance ./cmd/server

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash # For entrypoint script
WORKDIR /root/
COPY --from=builder /app/xyz-multifinance .
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
EXPOSE 3009
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["./xyz-multifinance"]