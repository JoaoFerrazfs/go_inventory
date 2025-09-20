FROM golang:1.24

WORKDIR /app

## Install air
RUN go install github.com/air-verse/air@latest

## Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH=$PATH:/go/bin

EXPOSE 3000

# baixa dependências e roda Air
CMD ["sh", "-c", "go mod download && air -c .air.toml", ]
