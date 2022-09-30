package main

import (
	"fmt"
	"html/template"
	"tpcmethod/http"
)

func main() {
	r := http.Router()

	tmp, err := template.ParseGlob("templates/*")
	if err != nil {
		fmt.Printf("Error al cargar el template")
	}

	r.Route("GET", "/", func(req *http.Request, res *http.Response) {
		//res.SendData("Hi")
		res.SendRender(tmp, "index.html")

	})

	r.Route("GET", "/ping", func(req *http.Request, res *http.Response) {
		res.SendJson("{\"ping\":\"pong\"}")
	})

	http.Serve(":8080", r)
}
