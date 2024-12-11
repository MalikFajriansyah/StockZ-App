# Dockerfile
default:
    # Menggunakan base image Ubuntu
    FROM ubuntu:22.04
    
    # Menentukan direktori kerja
    WORKDIR /app
    
    # Install dependensi yang diperlukan
    RUN apt-get update && apt-get install -y \
        curl \
        git \
        build-essential \
        wget && rm -rf /var/lib/apt/lists/*
    
    # Install Go
    RUN wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz && \
        tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz && \
        rm go1.21.1.linux-amd64.tar.gz
    
    # Menambahkan Go ke PATH
    ENV PATH "/usr/local/go/bin:$PATH"
    
    # Menyalin kode aplikasi ke dalam image
    COPY . .
    
    # Mendownload dependensi Go
    RUN go mod download
    
    # Membuild aplikasi Go
    RUN go build -o main .
    
    # Mendefinisikan port
    EXPOSE 8080
    
    # Menjalankan aplikasi
    CMD ["./main"]