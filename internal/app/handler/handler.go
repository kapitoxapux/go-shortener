package handler

import (
	"compress/gzip"

	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"

	"myapp/internal/app/config"

	"myapp/internal/app/storage"
	"net/http"
	"strings"

	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*chi.Mux
}

type JSONShorter struct {
	URL string `json:"url"`
}

type JSONBatcher struct {
	UrlID   string `json:"correlation_id"`
	LongURL string `json:"original_url"`
}

type JSONResultBatcher struct {
	UrlID    string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type JSONObject struct {
	Short_url    string `json:"short_url"`
	Original_url string `json:"original_url"`
}

func (w gzipWriter) Write(b []byte) (int, error) {

	return w.Writer.Write(b)
}

func SetUserCookie(req *http.Request, sign []byte) *http.Cookie {
	expiration := time.Now().Add(60 * time.Second)

	return &http.Cookie{
		Name:    "user_id",
		Value:   SetCookieToken(sign),
		Path:    req.URL.Path,
		Expires: expiration,
	}
}

func GetSignerCheck(sign []byte, cookie string) bool {
	return hmac.Equal(sign, TokenCheck(cookie))
}

func TokenCheck(cookie string) []byte {

	// 1) получите ключ из password, используя sha256.Sum256
	key := sha256.Sum256(config.Secretkey)

	// 2) создайте aesblock и aesgcm
	// будем использовать AES256, создав ключ длиной 32 байта
	aesblock, _ := aes.NewCipher(key[:])
	aesgcm, _ := cipher.NewGCM(aesblock)

	// 3) получите вектор инициализации aesgcm.NonceSize() байт с конца ключа
	nonceSize := aesgcm.NonceSize()
	nonce := key[len(key)-nonceSize:]

	// 4) декодируйте сообщение msg в двоичный формат
	dst, _ := hex.DecodeString(cookie)

	// 5) расшифруйте и выведите данные
	src, _ := aesgcm.Open(nil, nonce, dst, nil) // расшифровываем

	return src
}

func SetCookieToken(sign []byte) string {

	key := sha256.Sum256(config.Secretkey) // ключ шифрования

	aesblock, _ := aes.NewCipher(key[:32])
	aesgcm, _ := cipher.NewGCM(aesblock)

	// создаём вектор инициализации
	nonceSize := aesgcm.NonceSize()
	nonce := key[len(key)-nonceSize:]

	dst := aesgcm.Seal(nil, nonce, sign, nil) // симметрично зашифровываем

	return hex.EncodeToString(dst)
}

func GzipMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		defer gzw.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzw}, r)
	})
}

func SetShortAction(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed for this route!", http.StatusMethodNotAllowed)

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
			http.Error(res, err.Error(), http.StatusBadRequest)

			return
		}

		reader = gzr
		defer gzr.Close()
	} else {
		reader = req.Body
	}

	defer req.Body.Close()

	if req.ContentLength > 0 {
		b, err := io.ReadAll(reader)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)

			return
		}

		short, duplicate := storage.SetShort(string(b))

		cookie, _ := req.Cookie("user_id")
		if cookie == nil {
			http.SetCookie(res, SetUserCookie(req, short.Signer.Sign))
		} else {
			if !GetSignerCheck(short.Signer.Sign, cookie.Value) {
				http.SetCookie(res, SetUserCookie(req, short.Signer.Sign))
			}

		}

		if duplicate {
			res.WriteHeader(http.StatusConflict)
		} else {
			res.WriteHeader(http.StatusCreated)
		}

		res.Write([]byte(short.ShortURL))
	} else {
		http.Error(res, "Empty body!", http.StatusBadRequest)
	}
}

func GetShortAction(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed for this route", http.StatusMethodNotAllowed)

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
		http.Error(res, "Only POST requests are allowed for this route!", http.StatusMethodNotAllowed)

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

	if req.ContentLength > 0 {
		b, err := io.ReadAll(reader)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)

			return
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.Header().Add("Accept", "application/json")

		j := new(JSONShorter)
		if err := json.Unmarshal(b, &j); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}

		short, duplicate := storage.SetShort(j.URL)

		cookie, _ := req.Cookie("user_id")
		if cookie == nil {
			http.SetCookie(res, SetUserCookie(req, short.Signer.Sign))
		} else {
			if !GetSignerCheck(short.Signer.Sign, cookie.Value) {
				http.SetCookie(res, SetUserCookie(req, short.Signer.Sign))
			}

		}

		if duplicate {
			res.WriteHeader(http.StatusConflict)
		} else {
			res.WriteHeader(http.StatusCreated)
		}

		res.Write([]byte(`{"result":"` + short.ShortURL + `"}`))
	} else {
		http.Error(res, "Empty body!", http.StatusBadRequest)
	}
}

func GetUserURLAction(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed for this route", http.StatusMethodNotAllowed)

		return
	}

	if req.URL.Path != "/api/user/urls" {
		http.Error(res, "Wrong route!", http.StatusNotFound)

		return
	}

	list := []JSONObject{}
	cookie, _ := req.Cookie("user_id")
	if cookie == nil {
		http.Error(res, "No content!", http.StatusNoContent)

		return
	} else {
		res.Header().Set("Content-Type", "application/json; charset=utf-8")

		// пройтись по всем записям и забрать нужные объекты используя куки юзера
		for _, short := range storage.GetFullList() {
			if GetSignerCheck(short.Signer.Sign, cookie.Value) {
				obj := JSONObject{}
				obj.Short_url = short.ShortURL
				obj.Original_url = short.LongURL
				list = append(list, obj)
			}

		}
	}

	if len(list) == 0 {
		http.Error(res, "No content!", http.StatusNoContent)

		return
	} else {
		p, _ := json.Marshal(list)
		res.Write([]byte(p))
	}
}

func GetPingAction(res http.ResponseWriter, req *http.Request) {
	if status, err := storage.ConnectionDBCheck(); status != http.StatusOK {
		http.Error(res, err, http.StatusInternalServerError)

		return
	} else {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("OK"))

		return
	}

}

func GetBatchAction(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed for this route!", http.StatusMethodNotAllowed)

		return
	}

	if req.URL.Path != "/api/shorten/batch" {
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

	if req.ContentLength > 0 {
		b, err := io.ReadAll(reader)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)

			return
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.Header().Add("Accept", "application/json")

		list := []JSONBatcher{}

		// множество обьектов в теле нужен цикл
		if err := json.Unmarshal(b, &list); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}

		resultsObj := []JSONResultBatcher{}

		for i, obj := range list {
			short, _ := storage.SetShort(obj.LongURL)

			if i == 1 {
				cookie, _ := req.Cookie("user_id")
				if cookie == nil {
					http.SetCookie(res, SetUserCookie(req, short.Signer.Sign))
				} else {
					if !GetSignerCheck(short.Signer.Sign, cookie.Value) {
						http.SetCookie(res, SetUserCookie(req, short.Signer.Sign))
					}

				}
			}

			resultBatcher := new(JSONResultBatcher)
			resultBatcher.UrlID = obj.UrlID
			resultBatcher.ShortURL = short.ShortURL

			resultsObj = append(resultsObj, *resultBatcher)
		}

		res.WriteHeader(http.StatusCreated)
		if len(resultsObj) == 0 {
			http.Error(res, "No content!", http.StatusNoContent)

			return
		} else {
			p, _ := json.Marshal(resultsObj)
			res.Write([]byte(p))
		}
	} else {
		http.Error(res, "Empty body!", http.StatusBadRequest)
	}
}

func NewRoutes() *Handler {
	mux := &Handler{
		Mux: chi.NewMux(),
	}

	mux.Post("/", SetShortAction)
	mux.Get("/{`\\w+$`}", GetShortAction)
	mux.Post("/api/shorten", GetJSONShortAction)
	mux.Get("/api/user/urls", GetUserURLAction)
	mux.Get("/ping", GetPingAction)
	mux.Post("/api/shorten/batch", GetBatchAction)

	return mux
}
