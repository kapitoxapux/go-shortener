package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

var paths = map[string]*Shorter{}

type Shorter struct {
	id       string
	longUrl  string
	shortUrl string
}

type Handler struct {
	*chi.Mux
}

func NewHandler() *Handler {
	chi := &Handler{
		Mux: chi.NewMux(),
	}

	chi.Get("/{`\\w+$`}", chi.NewAction())

	chi.Post("/", chi.NewAction())

	return chi
}

func (h *Handler) NewAction() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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
}

func shortener(url string) string {
	plainText := []byte(url)
	sha256Hash := sha256.Sum256(plainText)

	return hex.EncodeToString(sha256Hash[:])
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

func NewRoutes() *Handler {
	return NewHandler()
}
