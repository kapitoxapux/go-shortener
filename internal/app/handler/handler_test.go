package handler

import (
	"io"
	"myapp/internal/app/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCustomAction(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
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
			http.Error(res, "Only GET and POST requests are allowed!", http.StatusNotFound)

			return
		}

	}

}

func TestEndpoints_Handle(t *testing.T) {

	forTest := storage.SetShort("http://localhost:8080/some_text_to_test_2")

	type want struct {
		contentType string
		statusCode  int
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
				statusCode:  400,
			},
			pattern: "/",
		},
		{
			name:   "simple test #2",
			method: "POST",
			body:   "fullUrl",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
			},
			pattern: "/",
		},
		{
			name:   "simple test #3",
			method: "GET",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  307,
			},
			pattern: forTest.ShortURL,
		},
		{
			name:   "simple test #4",
			method: "POST",
			body:   "fullUrl",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
			pattern: "/shortUrl",
		},
		{
			name:   "simple test #5",
			method: "PUT",
			body:   "fullUrl",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
			pattern: "/shortUrl",
		},
		{
			name:   "simple test #6",
			method: "PUT",
			body:   "fullUrl",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
			pattern: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(tt.method, tt.pattern, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(testCustomAction)
			h(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if request.Method == "POST" {
				_, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
			}

		})
	}
}
