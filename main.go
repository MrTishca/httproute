package main

import (
	"tpcmethod/http"
)

func main() {
	r := http.Router()
	r.Route("GET", "/", func(req http.Request, res http.Response) {
		res.SendData("Hi")
	})

	r.Route("GET", "/ping", func(req http.Request, res http.Response) {
		res.SendJson("{\"ping\":\"pong\"}")
	})

	http.Serve(":8080", r)
}
