package handler

import (
	"io"

	"myapp/internal/app/storage"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*chi.Mux
}

func PostAction(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed for this route!", http.StatusNotFound)

		return
	}

	if req.URL.Path != "/" {
		http.Error(res, "Wrong route!", http.StatusNotFound)

		return
	}

	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusCreated)

	short := storage.SetShort(string(b))

	res.Write([]byte(short.ShortURL))

}

func GetAction(res http.ResponseWriter, req *http.Request) {

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

func NewRoutes() *Handler {

	chi := &Handler{
		Mux: chi.NewMux(),
	}

	chi.Get("/{`\\w+$`}", GetAction)
	chi.Post("/", PostAction)

	return chi
}
