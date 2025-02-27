FROM golang:1.21-alpine

WORKDIR /app

# Adicionar git para poder baixar as dependências
RUN apk add --no-cache git

COPY go.mod ./
# Remover a cópia do go.sum já que ele pode não existir ainda
# COPY go.sum ./

# Baixar todas as dependências
RUN go mod download
RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/server

EXPOSE 8080

CMD ["./main"] 