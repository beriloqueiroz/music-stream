# Usar a versão específica do Go
FROM golang:1.22.2-alpine

WORKDIR /app

# Instalar dependências necessárias
RUN apk add --no-cache git protobuf-dev

# Copiar os arquivos de definição de módulo
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar o resto do código
COPY . .

# Criar diretório para armazenamento de músicas
RUN mkdir -p storage/music

# Compilar a aplicação
RUN go build -o main ./cmd/server

# Expor as portas HTTP e gRPC
EXPOSE 8080 50051

CMD ["./main"]