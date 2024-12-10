# Use a imagem oficial do Go
FROM golang:1.18

# Define o diretório de trabalho
WORKDIR /app

# Copie o go.mod e o go.sum
COPY go.mod go.sum ./

# Instale as dependências
RUN go mod download

# Copie o código-fonte
COPY . .

# Compile o aplicativo
RUN go build -o main .

# Exponha a porta
EXPOSE 8080

# Execute o binário
CMD ["./main"]
