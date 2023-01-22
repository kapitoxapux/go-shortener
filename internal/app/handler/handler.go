package handler

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"myapp/internal/app/storage"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*chi.Mux
}

type JSONShorter struct {
	URL string `json:"url"`
}

var j JSONShorter

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gzw.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzw}, r)
	})
}

func SetShortAction(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed for this route!", http.StatusNotFound)

		return
	}

	if req.URL.Path != "/" {
		http.Error(res, "Wrong route!", http.StatusNotFound)

		return
	}

	var reader io.Reader

	if req.Header.Get(`Content-Encoding`) == `gzip` {
		gzr, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gzr
		defer gzr.Close()
	} else {
		reader = req.Body
	}

	defer req.Body.Close()

	b, err := io.ReadAll(reader)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusCreated)

	short := storage.SetShort(string(b))

	res.Write([]byte(short.ShortURL))

}

func GetShortAction(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed for this route", http.StatusNotFound)

		return
	}

	part := req.URL.Path
	formated := strings.Replace(part, "/", "", -1)

	sh := storage.GetShort(formated)
	if sh == "" {
		http.Error(res, "Url not founded!", http.StatusBadRequest)

		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set("Location", storage.GetFullURL(formated))
	res.WriteHeader(http.StatusTemporaryRedirect)

}

func GetJSONShortAction(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed for this route!", http.StatusNotFound)

		return
	}

	if req.URL.Path != "/api/shorten" {
		http.Error(res, "Wrong route!", http.StatusNotFound)

		return
	}

	defer req.Body.Close()

	var reader io.Reader

	if req.Header.Get(`Content-Encoding`) == `gzip` {
		gzr, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gzr
		defer gzr.Close()
	} else {
		reader = req.Body
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Add("Accept", "application/json")

	if err := json.Unmarshal(b, &j); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	short := storage.SetShort(j.URL)

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(`{"result":"` + short.ShortURL + `"}`))
}

func NewRoutes() *Handler {
	mux := &Handler{
		Mux: chi.NewMux(),
	}

	mux.Post("/", SetShortAction)
	mux.Get("/{`\\w+$`}", GetShortAction)
	mux.Post("/api/shorten", GetJSONShortAction)

	return mux
}
