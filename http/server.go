package http

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Request struct {
	Method  string
	Route   string
	Headers map[string]string
	Host    string
	Secure  string
	Body    []string
}

type Response struct {
	Method  string
	Route   string
	headers map[string]string
	Host    string
}

type handler struct {
	path    string
	handler func(Request, Response)
}

type Router struct {
	list []handler
}

var router Router = Router{}

func (r Router) Get(path string, handler func(Request, Response)) {

	if len(r.list) == 0 {
		// r.list[0] = handler {
		// 	path   : path,
		// 	handler :handler,
		// }
	} else {
		r.list[len(r.list)+1] = struct {
			path    string
			handler func(Request, Response)
		}{path, handler}
	}
}

func build_request(conn net.Conn) Request {
	h, b := split_request(conn)
	req := split_headers(h)
	req.Body = b
	return req
}

func start_router(conn net.Conn) {
	//codigo chido aquis
	defer conn.Close()
	req := build_request(conn)
	respond(conn, req)
}

func split_request(conn net.Conn) ([]string, []string) {
	scanner := bufio.NewScanner(conn)
	var h []string
	var b []string
	f := true

	for scanner.Scan() {
		ln := scanner.Text()
		if ln == "" {
			f = false
		}
		if f == true {
			h = append(h, ln)
		}
		if f == false {
			b = append(b, ln)
		}
		if f == false && ln == "" {
			break
		}
	}
	return h, b
}

func split_headers(hd []string) Request {
	req := Request{}
	req.Headers = make(map[string]string)
	for i, h := range hd {
		fmt.Println(h)
		if i == 0 {
			spl := strings.Split(h, " ")
			req.Method = spl[0]
			req.Route = spl[1]
			req.Secure = spl[2]
		} else {
			spl := strings.Split(h, ": ")
			if spl[0] == "Host" {
				req.Host = spl[1]
			} else {
				req.Headers[spl[0]] = spl[1]
			}

		}
	}
	return req
}

func respond(conn net.Conn, req Request) {

	body := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title></title></head><body><strong>` + req.Host + req.Route + `</strong> <br> <span>conection from:` + req.Headers["User-Agent"] + `</span></body></html>`

	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, body)
}

func Serve(port string, r Router) {
	router = r
	serve(port)
}

func serve(port string) {
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err.Error())
			continue
		}
		go start_router(conn)
	}
}
