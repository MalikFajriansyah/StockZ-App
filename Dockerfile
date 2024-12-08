# Gunakan base image untuk Golang
FROM golang:1.23.4 AS builder

# Set working directory
WORKDIR /app

# Copy semua file ke dalam container
COPY . .

# Download dan install dependencies
RUN go mod tidy

# Build aplikasi
RUN go build -o main .

# Image untuk produksi
FROM debian:bullseye

# Install CA certificates (penting untuk Firebase)
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

# Set working directory
WORKDIR /app

# Copy binary dari tahap builder
COPY --from=builder /app/main .

# Copy file konfigurasi, jika ada
COPY ServiceAccountKey.json /app/ServiceAccountKey.json

# Expose port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
