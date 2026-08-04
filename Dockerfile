FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/doq ./cmd/doq.go

FROM alpine:3.24.1

WORKDIR /app

COPY --from=builder /app/doq /app/doq

EXPOSE 8080

USER nobody:nobody

ENTRYPOINT ["/app/doq"]

