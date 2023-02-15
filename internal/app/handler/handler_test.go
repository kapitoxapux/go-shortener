package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myapp/internal/app/config"
	"myapp/internal/app/repository"
	"myapp/internal/app/service"
	"myapp/internal/app/storage"
)

var forTest *service.Shorter
var repo repository.Repository

func GetDB() service.Storage {

	if status, _ := ConnectionDBCheck(); status == http.StatusOK {
		return storage.NewDB()
	}

	if pathStorage := config.GetConfigPath(); pathStorage != "" {
		return storage.NewFileDB()
	}

	return storage.NewInMemDB()
}

func testCustomAction(res http.ResponseWriter, req *http.Request) {
	if status, _ := ConnectionDBCheck(); status == http.StatusOK {
		repo = repository.NewRepository(config.GetStorageDB())
	} else {
		repo = nil
	}

	switch req.URL.Path {
	case "/":
		if req.Method != http.MethodPost {
			http.Error(res, "Wrong route!", http.StatusMethodNotAllowed)

			return
		}

		if req.URL.Path != "/" {
			http.Error(res, "Wrong route!", http.StatusNotFound)

			return
		}

		defer req.Body.Close()
		_, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)

			return
		}

		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusCreated)

		res.Write([]byte(forTest.ShortURL))

	case "/api/shorten":
		if req.Method != http.MethodPost {
			http.Error(res, "Wrong route!", http.StatusMethodNotAllowed)

			return
		}

		defer req.Body.Close()
		_, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)

			return
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.Header().Add("Accept", "application/json")

		res.WriteHeader(http.StatusCreated)

		res.Write([]byte(`{"result":"` + forTest.ShortURL + `"}`))

	default:
		if req.Method != http.MethodGet {
			http.Error(res, "Wrong route!", http.StatusMethodNotAllowed)

			return
		}

		part := req.URL.Path
		formated := strings.Replace(part, "/", "", -1)

		db := GetDB()
		service := service.NewService(db)

		sh := service.Storage.GetShort(formated)
		if sh == "" {
			http.Error(res, "Url not founded!", http.StatusBadRequest)

			return
		}

		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.Header().Set("Location", service.Storage.GetFullURL(formated))
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func TestEndpoints_Handle(t *testing.T) {
	db := GetDB()
	service := service.NewService(db)

	if status, _ := ConnectionDBCheck(); status == http.StatusOK {
		repo = repository.NewRepository(config.GetStorageDB())
	} else {
		repo = nil
	}

	forTest, _ = service.Storage.SetShort("https://dev.to/nwneisen/writing-a-url-shortener-in-go-2ld6")

	type want struct {
		contentType string
		statusCode  int
		bodyContent string
	}

	tests := []struct {
		name    string
		method  string
		body    string
		pattern string
		want    want
	}{
		{
			name:   "simple test #1",
			method: "GET",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  405,
				bodyContent: "Wrong route!\n",
			},
			pattern: config.GetConfigBase() + "/",
		},
		{
			name:   "simple test #2",
			method: "POST",
			body:   forTest.LongURL,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				bodyContent: forTest.ShortURL,
			},
			pattern: config.GetConfigBase() + "/",
		},
		{
			name:   "simple test #3",
			method: "GET",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				bodyContent: "",
			},
			pattern: forTest.ShortURL,
		},
		{
			name:   "simple test #4",
			method: "POST",
			body:   forTest.LongURL,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  405,
				bodyContent: "Wrong route!\n",
			},
			pattern: forTest.ShortURL,
		},
		{
			name:   "simple test #5",
			method: "PUT",
			body:   forTest.LongURL,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  405,
				bodyContent: "",
			},
			pattern: forTest.ShortURL,
		},
		{
			name:   "simple test #6",
			method: "PUT",
			body:   forTest.LongURL,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  405,
				bodyContent: "Wrong route!\n",
			},
			pattern: config.GetConfigBase() + "/",
		},
		{
			name:   "simple test #7",
			method: "POST",
			body:   `{"url":"` + forTest.LongURL + `"}`,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  405,
				bodyContent: "Wrong route!\n",
			},
			pattern: "/api",
		},
		{
			name:   "simple test #8",
			method: "POST",
			body:   `{"url":"` + forTest.LongURL + `"}`,
			want: want{
				contentType: "application/json; charset=utf-8",
				statusCode:  201,
				bodyContent: `{"result":"` + forTest.ShortURL + `"}`,
			},
			pattern: config.GetConfigBase() + "/api/shorten",
		},
	}

	var request *http.Request
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.method == "POST" {
				request = httptest.NewRequest(tt.method, tt.pattern, bytes.NewBuffer([]byte(tt.body)))
			} else {
				request = httptest.NewRequest(tt.method, tt.pattern, nil)
			}

			w := httptest.NewRecorder()

			h := http.HandlerFunc(testCustomAction)

			h(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if request.Method == "POST" {
				defer result.Body.Close()
				link, err := io.ReadAll(result.Body)
				require.NoError(t, err)

				assert.Equal(t, tt.want.bodyContent, string(link))

			}

		})
	}
}
