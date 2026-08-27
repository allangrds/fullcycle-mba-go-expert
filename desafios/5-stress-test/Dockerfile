# builder: contém o toolchain Go completo + código-fonte, usado só para compilar o binário.
FROM golang:1.26-alpine AS builder

WORKDIR /app

# gcc/musl-dev: necessários para compilar com "-race" (cgo) ao rodar os
# testes via "docker compose run --rm test". Não engordam a imagem final:
# o stage abaixo só copia o binário compilado, sem essas libs de build.
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /stress-test .

# imagem final: só o binário, sem toolchain Go, bem menor que a de build.
FROM alpine:latest

COPY --from=builder /stress-test /stress-test

# ENTRYPOINT (e não CMD): os argumentos passados em
# "docker run <imagem> --url=... --requests=... --concurrency=..." são
# anexados ao binário, em vez de substituir o comando inteiro do container.
ENTRYPOINT ["/stress-test"]
