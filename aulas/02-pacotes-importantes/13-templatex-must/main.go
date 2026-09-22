package main

import (
	"os"
	"text/template"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

func main() {
	curso := Curso{
		Nome:         "Go",
		CargaHoraria: 40,
	}

	// Must é um método do pacote template que cria um novo template e trata erros de forma mais simples
	// Se ocorrer um erro, o programa irá panicar, o que é útil para evitar
	// verificações de erro repetitivas e deixar o código mais limpo.
	templat := template.Must(template.New("curso_template").Parse("Curso: {{.Nome}} | Carga Horária: {{.CargaHoraria}}"))
	// templat := template.New("curso_template")
	// templat, _ = templat.Parse("Curso: {{.Nome}} | Carga Horária: {{.CargaHoraria}}")

	// O código abaixo executa o template e escreve a saída no os.Stdout
	err := templat.Execute(os.Stdout, curso)
	if err != nil {
		panic(err)
	}
}
