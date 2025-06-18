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

	templat := template.New("curso_template")
	templat, _ = templat.Parse("Curso: {{.Nome}} | Carga Horária: {{.CargaHoraria}}")

	// O código abaixo executa o template e escreve a saída no os.Stdout
	err := templat.Execute(os.Stdout, curso)
	if err != nil {
		panic(err)
	}
}
