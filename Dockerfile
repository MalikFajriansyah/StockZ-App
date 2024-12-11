# Gunakan image base Ubuntu dengan versi yang sesuai
FROM ubuntu:22.04

# Update dan upgrade paket sistem
RUN apt-get update && apt-get upgrade -y

# Install Go
RUN wget https://go.dev/dl/go1.20.linux-amd64.tar.gz \
    && tar -xvf go1.20.linux-amd64.tar.gz \
    && mv go /usr/local \
    && rm -rf go1.20.linux-amd64.tar.gz \
    && export PATH=$PATH:/usr/local/go/bin

# Install git untuk clone repository
RUN apt-get install git -y

# Buat direktori kerja dan pindah ke sana
WORKDIR /app

# Clone repository Go project Anda
COPY . .

# Install dependensi Go
RUN go mod download

# Build aplikasi Go
RUN go build -o main .

# Install PostgreSQL
RUN apt-get install postgresql postgresql-contrib -y

# Inisialisasi database PostgreSQL (sesuaikan dengan konfigurasi Anda)
RUN psql -c "CREATE DATABASE your_database_name;" \
    && psql -c "CREATE USER your_username WITH PASSWORD 'your_password';" \
    && psql -c "GRANT ALL PRIVILEGES ON DATABASE your_database_name TO your_username;"

# Copy file konfigurasi database (jika diperlukan)
COPY postgresql.conf /etc/postgresql/postgresql.conf

# Set environment variable untuk koneksi database
ENV DATABASE_URL=postgres://your_username:your_password@localhost:5432/your_database_name

# Expose port untuk koneksi PostgreSQL
EXPOSE 5432

# Jalankan aplikasi Go
CMD ["/app/main"]