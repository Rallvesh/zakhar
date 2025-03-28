FROM golang:alpine AS builder

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY . .
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/zakhar ./cmd/zakhar/




FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/zakhar /go/bin/zakhar
ENTRYPOINT ["/go/bin/zakhar"]
