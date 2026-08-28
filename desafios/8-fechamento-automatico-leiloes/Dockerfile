# builder: contém o toolchain Go completo + código-fonte. É reaproveitado
# pelo serviço "test" do docker-compose (--target builder) para rodar
# `go test ./...` sem precisar de uma segunda imagem.
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /auction ./cmd/auction

# imagem final: só o binário, bem menor que a de build.
FROM alpine:latest

WORKDIR /app

COPY --from=builder /auction .

EXPOSE 8080

CMD ["./auction"]
