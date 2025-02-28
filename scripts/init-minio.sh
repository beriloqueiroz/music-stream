#!/bin/sh

# Esperar o MinIO iniciar
sleep 5

# Instalar cliente MinIO
wget https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x mc

# Configurar cliente
./mc alias set myminio http://minio:9000 minioadmin minioadmin

# Criar bucket se não existir
./mc mb --ignore-existing myminio/music-bucket

# Limpar
rm mc 