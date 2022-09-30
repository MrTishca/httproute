package http

import (
	"bufio"
	"fmt"
	"html/template"
	templ "html/template"
	"log"
	"net"
	"strings"
)

type Request struct {
	Method  string
	Route   string
	Params  map[string]string
	Headers map[string]string
	Host    string
	Secure  string
	Body    []string
}

type Response struct {
	headers map[string]string
	connect net.Conn
}

type handler func(*Request, *Response)

type key struct {
	method string
	path   string
}

type router_ struct {
	list map[key]handler
}

func (r router_) Route(method string, path string, hand func(*Request, *Response)) {
	method = strings.ToUpper(method)
	r.list[key{method, path}] = hand
}

func Router() router_ {
	r := router_{
		make(map[key]handler),
	}
	hand := func(req *Request, res *Response) {
		res.SendData("Welcome to 404")
	}
	r.list[key{"GET", "/404"}] = hand
	return r
}

func build_request(conn net.Conn) Request {
	h, b := split_request(conn)
	req := split_headers(h)
	params := strings.Split(req.Route, "?")
	if len(params) > 1 {
		req.Route = params[0]
		req.Params = split_params(params[1])
	}
	req.Body = b
	return req
}

func mutltiplexor(r router_, req Request, res Response) {
	k := key{req.Method, req.Route}
	if _, ok := r.list[k]; ok {
		r.list[key{req.Method, req.Route}](&req, &res)
	}

	// fmt.Printf("%v", pos_route)
	// if len(pos_route) == 0 {
	// 	fmt.Println("Join in 404")
	// 	r.list[key{"GET", "/404"}](req, res)
	// } else if len(pos_route) == 1 {
	// 	fmt.Println("Join in single route")
	// 	r.list[key{req.Method, req.Route}](req, res)
	// } else {
	// 	//Multiple case for  params:id feature!!

	// }
}

func split_params(data string) map[string]string {
	m := make(map[string]string)
	spli := strings.Split(data, "&")
	fmt.Printf("%v", spli)
	for _, s := range spli {
		d := strings.Split(s, "=")
		m[d[0]] = d[1]
	}
	return m
}

func start_router(conn net.Conn, r router_) {
	defer conn.Close()
	req := build_request(conn)
	res := Response{make(map[string]string), conn}
	mutltiplexor(r, req, res)
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

func (r Response) SendData(data string) {
	fmt.Fprint(r.connect, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(r.connect, "Content-Length: %d\r\n", len(data))
	fmt.Fprint(r.connect, "Content-Type: text/html\r\n")
	fmt.Fprint(r.connect, "\r\n")
	fmt.Fprint(r.connect, data)
}
func (r Response) SendJson(data string) {
	fmt.Fprint(r.connect, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(r.connect, "Content-Length: %d\r\n", len(data))
	fmt.Fprint(r.connect, "Content-Type: application/json\r\n")
	fmt.Fprint(r.connect, "\r\n")
	fmt.Fprint(r.connect, data)
}
func (r Response) SendRender(tmp *templ.Template, t string) {
	fmt.Fprint(r.connect, "HTTP/1.1 200 OK\r\n")
	fmt.Fprint(r.connect, "Content-Type: text/html; charset=utf-8\r\n")
	fmt.Fprint(r.connect, "Access-Control-Allow-Origin: *\r\n")
	fmt.Fprint(r.connect, "Connection: Keep-Alive\r\n")
	fmt.Fprint(r.connect, "Keep-Alive: timeout=5, max=997\r\n")
	fmt.Fprint(r.connect, "Server: Apache\r\n")
	fmt.Fprint(r.connect, "\r\n")

	tpl := template.Must(tmp, nil)
	err := tpl.ExecuteTemplate(r.connect, t, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func Serve(port string, r router_) {
	serve(port, r)
}

func serve(port string, r router_) {
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
		go start_router(conn, r)
	}
}
