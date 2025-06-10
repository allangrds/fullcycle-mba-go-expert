# Fullcycle | MBA | GO Expert

## Comandos Go CLI

- `go env`: show Go configured env variables
  - Example: `GOOS` variable

- `go run`: Compila e executa código Go
  - Exemplo: `go run main.go`

- `go mod`: Gerenciamento do uso de módulos no projeto
  - Exemplo: `go mod init <name>`

- `go mod tidy` - baixa as dependências e atualiza arquivo `go.mod`

- `go get`: Baixa packages importados e suas dependências
  - Exemplo `go get github.com/google/uuid`

- `go build`: Faz o build do arquivo Go
  - Exemplo `go build`, se já tiver criado um módulo através do `go mod init <name>`
  - Exemplo `go build main.go`
  - Exemplo `GOOS=windows go build main.go`
  - Exemplo `GOOS=linux go build main.go`
  - https://www.digitalocean.com/community/tutorials/building-go-applications-for-different-operating-systems-and-architectures
