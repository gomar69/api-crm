# Gunakan base image golang
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod dan go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary
RUN go build -o api-crm ./cmd

# Stage production
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/api-crm .

EXPOSE 8080

CMD ["./api-crm"]
