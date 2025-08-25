# Build stage
FROM golang:1.24.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api-crm .

# Run stage
FROM golang:1.24.5
WORKDIR /app

COPY --from=builder /app/api-crm .
COPY .env .

EXPOSE 8080
CMD ["./api-crm"]
