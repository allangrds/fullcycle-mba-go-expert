package main

import (
	"os"
	"text/template"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

type Cursos []Curso

func main() {
	cursos := Cursos{
		{
			Nome:         "Go",
			CargaHoraria: 40,
		},
		{
			Nome:         "Python",
			CargaHoraria: 30,
		},
	}
	templates := []string{
		"header.html",
		"content.html",
		"footer.html",
	}

	templat := template.Must(template.New("content.html").ParseFiles(templates...))
	err := templat.Execute(os.Stdout, cursos)
	if err != nil {
		panic(err)
	}
}
