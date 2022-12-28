package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getFullUrl(t *testing.T) {

	*paths["text_to_test_2"] = Shorter{
		id:       "text_to_test_2",
		longUrl:  "http://localhost:8080/some_text_to_test_2",
		shortUrl: "http://localhost:8080/text_to_test_2",
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "catch error",
			id:      "some_text_to_test_1",
			wantErr: true,
		},
		{
			name:    "get Full Url",
			id:      "text_to_test_2",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFullUrl(tt.id); (got != "") != tt.wantErr {
				t.Errorf("getFullUrl() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func Test_getShort(t *testing.T) {

	*paths["text_to_test_2"] = Shorter{
		id:       "text_to_test_2",
		longUrl:  "http://localhost:8080/some_text_to_test_2",
		shortUrl: "http://localhost:8080/text_to_test_2",
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "catch error",
			id:      "some_text_to_test_1",
			wantErr: true,
		},
		{
			name:    "get Short Url",
			id:      "text_to_test_2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getShort(tt.id); (got != "") != tt.wantErr {
				t.Errorf("getShort() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func Test_setShort(t *testing.T) {

	toCheck := setShort("http://localhost:8080/some_text_to_test_2")

	tests := []struct {
		name string
		link string
		want string
	}{
		{
			name: "new Shorter",
			link: "http://localhost:8080/some_text_to_test_2",
			want: toCheck.id,
		},
		{
			name: "catch error",
			link: "http://localhost:8080/some_text_to_test_1",
			want: toCheck.id,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setShort(tt.link); got.id != tt.want {
				t.Errorf("setShort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_CustomAction(t *testing.T) {
	type args struct {
		req *Request
	}
	tests := []struct {
		name string
		res  *Response
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.res.CustomAction(tt.args.req)
		})
	}
}

func testCustomAction() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "POST":
			b, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(res, err.Error(), http.StatusNotFound)

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

func TestEndpoints_Handle(t *testing.T) {

	// type args struct {
	// 	pattern string
	// 	handler CustomHandler
	// }

	type want struct {
		contentType string
		statusCode  int
	}

	// var h

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
			body:   "",
			want: want{
				contentType: "text/plane",
				statusCode:  404,
			},
			pattern: "/",
		},
		{
			name:   "simple test #2",
			method: "POST",
			body:   "fullUrl",
			want: want{
				contentType: "text/plane",
				statusCode:  201,
			},
			pattern: "/",
		},
		{
			name:   "simple test #3",
			method: "GET",
			body:   "",
			want: want{
				contentType: "text/plane",
				statusCode:  307,
			},
			pattern: "/shortUrl",
		},
		{
			name:   "simple test #4",
			method: "POST",
			body:   "fullUrl",
			want: want{
				contentType: "text/plane",
				statusCode:  404,
			},
			pattern: "/shortUrl",
		},
		{
			name:   "simple test #5",
			method: "PUT",
			body:   "fullUrl",
			want: want{
				contentType: "text/plane",
				statusCode:  404,
			},
			pattern: "/shortUrl",
		},
		{
			name:   "simple test #6",
			method: "PUT",
			body:   "fullUrl",
			want: want{
				contentType: "text/plane",
				statusCode:  404,
			},
			pattern: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(tt.method, tt.pattern, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(testCustomAction())
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			// w := httptest.NewRecorder()
			// routes.Handle(tt.pattern, func(w *Response, request *http.Request) {
			// 	h := http.HandlerFunc(w.CustomAction(httptest.NewRequest(tt.method, tt.pattern, nil)))
			// })

		})
	}
}
