# Stage 1: Build Stage
FROM golang:1.23.4 AS builder

# Set working directory
WORKDIR /app

# Copy semua file ke dalam container
COPY . .

# Download dan install dependencies
RUN go mod tidy

# Build aplikasi
RUN go build -o main .

# Stage 2: Production Stage
FROM ubuntu:20.04

# Set non-interactive mode untuk apt
ENV DEBIAN_FRONTEND=noninteractive

# Update repository dan install dependencies untuk GLIBC
RUN sed -i 's/http:\/\/archive.ubuntu.com/http:\/\/mirror.kakao.com/' /etc/apt/sources.list \
    && apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    build-essential \
    && wget http://ftp.gnu.org/gnu/libc/glibc-2.32.tar.gz \
    && tar -xvzf glibc-2.32.tar.gz \
    && cd glibc-2.32 \
    && mkdir build \
    && cd build \
    && ../configure --prefix=/opt/glibc-2.32 \
    && make -j$(nproc) \
    && make install \
    && cd ../.. \
    && rm -rf glibc-2.32 glibc-2.32.tar.gz \
    && apt-get remove --purge -y wget build-essential \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set GLIBC baru sebagai default
ENV LD_LIBRARY_PATH=/opt/glibc-2.32/lib:$LD_LIBRARY_PATH

# Set working directory
WORKDIR /app

# Copy binary dari tahap builder
COPY --from=builder /app/main .

# Copy file konfigurasi Firebase
COPY ServiceAccountKey.json /app/ServiceAccountKey.json

# Expose port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
