FROM golang:1.22.3 as builder

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos go.mod e go.sum para o diretório de trabalho
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia os arquivos .go para o diretório de trabalho
COPY . .

# Compila o binário
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o app ./cmd/bank/main.go

# Cria uma imagem alpine
FROM alpine:latest

# Define o diretório de trabalho
WORKDIR /root/

# Copia o binário do builder para o diretório de trabalho
COPY --from=builder /app/app .

# Create a .env file
RUN touch .env

# Define o comando de execução
CMD ["./app"]