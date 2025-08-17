FROM golang:1.24

WORKDIR /app

RUN go install github.com/air-verse/air@latest

EXPOSE 3000

# baixa dependências e roda Air
CMD ["sh", "-c", "go mod download && air -c .air.toml"]
