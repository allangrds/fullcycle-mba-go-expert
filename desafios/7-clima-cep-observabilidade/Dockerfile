FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# estágio "build" só é executado quando CMD_PATH é passado via --build-arg
# (service-a e service-b no docker-compose.yaml); "builder" acima é reaproveitado
# puro pelo serviço "test", sem compilar nenhum binário.
FROM builder AS build
ARG CMD_PATH
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ${CMD_PATH}

FROM alpine:3.22

WORKDIR /app

COPY --from=build /server /server

CMD ["/server"]
