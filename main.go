package main

import (
	"fmt"
	"tpcmethod/http"
)

func main() {
	r := http.Router{}
	r.Get("/", func(r1 http.Request, r2 http.Response) {
		fmt.Printf("Hi")
	})

	http.Serve(":8080", r)
}
