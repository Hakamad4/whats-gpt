FROM golang:1.20

# Define o diretório de trabalho
WORKDIR /go/src/app

# Copia o código fonte para o diretório de trabalho
COPY . .

# Executa o script gomod.sh para gerar o arquivo go.mod
RUN chmod u+x gomod.sh && bash ./gomod.sh

# Compila o código para o sistema operacional Linux e arquitetura amd64
RUN env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/whatsgpt whatsgpt/main.go

# Define o comando de entrada do contêiner para executar a aplicação
CMD ["./bin/whatsgpt"]
