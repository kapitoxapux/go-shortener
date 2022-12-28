package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	// "myapp/pkg/handler"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var paths = map[string]*Shorter{}

type CustomHandler func(*Response, *Request)

type Shorter struct {
	id       string
	longUrl  string
	shortUrl string
}

type Route struct {
	Pattern *regexp.Regexp
	Handler CustomHandler
}

type Endpoints struct {
	Routes []Route
	// DefaultRoute CustomHandler
}

type Request struct {
	*http.Request
	Params []string
}

type Response struct {
	http.ResponseWriter
}

func (r *Endpoints) Handle(pattern string, handler CustomHandler) {
	re := regexp.MustCompile(pattern)
	route := Route{
		Pattern: re,
		Handler: handler,
	}

	r.Routes = append(r.Routes, route)
}

func EndpointsHandler() *Endpoints {
	points := &Endpoints{

		// DefaultRoute: func(resp *Response, req *Request) {
		// 	resp.CustomAction(req)
		// },

	}

	return points
}

func (p *Endpoints) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &Request{Request: r}
	resp := &Response{w}

	for _, rt := range p.Routes {
		if matches := rt.Pattern.FindStringSubmatch(r.URL.Path); len(matches) > 0 {
			rt.Handler(resp, req)

			return
		}
	}
}

func (res *Response) CustomAction(req *Request) {
	switch req.Method {
	case "POST":

		if req.URL.Path != "/" {
			http.Error(res, "Wrong route!", http.StatusBadRequest)

			return
		}

		b, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)

			return
		}

		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusCreated)

		short := setShort(string(b))

		res.Write([]byte(short.shortUrl))

	case "GET":
		part := req.URL.Path
		formated := strings.Replace(part, "/", "", -1)

		sh := getShort(formated)
		if sh == "" {
			http.Error(res, "Url not founded!", http.StatusBadRequest)

			return
		}

		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.Header().Set("Location", getFullUrl(formated))
		res.WriteHeader(http.StatusTemporaryRedirect)

	default:
		if req.Method != http.MethodGet {
			http.Error(res, "Only GET and POST requests are allowed!", http.StatusBadRequest)

			return
		}

	}
}

func shortener(url string) string {
	b := make([]byte, 9)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(url))]
	}

	return string(b)
}

func setShort(url string) *Shorter {

	id := shortener(url)

	shorter := new(Shorter)

	shorter.id = id
	shorter.shortUrl = "http://localhost:8080/" + id
	shorter.longUrl = url

	paths[shorter.id] = shorter

	return shorter
}

func getShort(id string) string {
	if paths[id] != nil {
		return paths[id].shortUrl
	}
	return ""
}

func getFullUrl(id string) string {
	if paths[id] != nil {
		return paths[id].longUrl
	}
	return ""
}

func main() {

	routes := EndpointsHandler()

	routes.Handle(`\w+$`, func(resp *Response, req *Request) {
		resp.CustomAction(req)
	})

	routes.Handle("/", func(resp *Response, req *Request) {
		resp.CustomAction(req)
	})

	// mux := http.NewServeMux()
	// mux.Handle("/", http.HandlerFunc(PostAction))
	// mux.Handle("/*", http.HandlerFunc(GetAction))

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: routes,
	}
	log.Fatal(server.ListenAndServe())

}
