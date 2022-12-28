package handler

import (
	"io"
	"myapp/pkg/storage"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

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

			short := storage.SetShort(string(b))

			res.Write([]byte(short.ShortURL))

		case "GET":
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

		default:
			if req.Method != http.MethodGet {
				http.Error(res, "Only GET and POST requests are allowed!", http.StatusBadRequest)

				return
			}

		}
	}
}

func NewRoutes() *Handler {
	return NewHandler()
}
