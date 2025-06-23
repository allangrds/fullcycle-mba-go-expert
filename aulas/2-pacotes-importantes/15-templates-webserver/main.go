package main

import (
	"net/http"
	"text/template"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

type Cursos []Curso

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templat := template.Must(template.New("template.html").ParseFiles("template.html"))
		err := templat.Execute(w, Cursos{
			{
				Nome:         "Go",
				CargaHoraria: 40,
			},
			{
				Nome:         "Python",
				CargaHoraria: 30,
			},
		})
		if err != nil {
			panic(err)
		}
	})
	http.ListenAndServe(":8282", nil)
}
