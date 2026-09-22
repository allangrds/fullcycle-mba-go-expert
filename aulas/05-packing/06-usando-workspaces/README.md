Ao invés de fazer o `go mod edit -replace` vamos usar o `go work init ./math ./system` na raíz do `6-usando-workspaces`

Só com isso já vai funcionar o `go run system/main.go`

O problema é que, se você tiver outros pacotes já publicados, ao rodar o `go mod tidy` não vai baixar por conta dos pacotes locais ainda não publicados. Quais são as soluções?

1. Publicar os pacotes(mas perde o sentido do workspace)
2. Executar `go mod tidy -e` pra ele ignorar os pacotes que não achou
