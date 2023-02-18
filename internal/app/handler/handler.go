package handler

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"myapp/internal/app/config"
	"myapp/internal/app/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {

	return &Handler{
		service: service,
	}
}

type JSONShorter struct {
	URL string `json:"url"`
}

type JSONObject struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type JSONBatcher struct {
	URLID   string `json:"correlation_id"`
	LongURL string `json:"original_url"`
}

type JSONResultBatcher struct {
	URLID    string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {

	return w.Writer.Write(b)
}

func SetUserCookie(req *http.Request, sign []byte) *http.Cookie {
	expiration := time.Now().Add(6000 * time.Second)

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
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)

					return
				}
				w.Header().Set("Content-Encoding", "gzip")
				w = gzipWriter{
					ResponseWriter: w,
					Writer:         gzw,
				}
				defer gzw.Close()
			}
			if r.Header.Get("Content-Encoding") == "gzip" {
				gzr, err := gzip.NewReader(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}
				r.Body = gzr
				defer gzr.Close()
			}
			next.ServeHTTP(w, r)
		},
	)
}

func ConnectionDBCheck() (int, string) {
	db, err := sql.Open("pgx", config.GetStorageDB())
	if err != nil {

		return 500, err.Error()
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {

		return 500, err.Error()
	}

	return 200, ""
}

func (h *Handler) SetShortAction(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed for this route!", http.StatusMethodNotAllowed)

		return
	}
	if req.URL.Path != "/" {
		http.Error(res, "Wrong route!", http.StatusNotFound)

		return
	}
	defer req.Body.Close()
	if req.ContentLength < 1 {
		http.Error(res, "Empty body!", http.StatusBadRequest)
	}
	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	short, duplicate := h.service.Storage.SetShort(string(b))
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
}

func (h *Handler) GetShortAction(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed for this route", http.StatusMethodNotAllowed)

		return
	}
	part := req.URL.Path
	formated := strings.Replace(part, "/", "", -1)
	sh := h.service.Storage.GetShort(formated)
	if sh == "" {
		http.Error(res, "Url not founded!", http.StatusBadRequest)

		return
	}
	if sh == "402" {
		res.WriteHeader(http.StatusGone)

		return
	}
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set("Location", h.service.Storage.GetFullURL(formated))
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) GetJSONShortAction(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed for this route!", http.StatusMethodNotAllowed)

		return
	}
	if req.URL.Path != "/api/shorten" {
		http.Error(res, "Wrong route!", http.StatusNotFound)

		return
	}
	defer req.Body.Close()
	if req.ContentLength < 1 {
		http.Error(res, "Empty body!", http.StatusBadRequest)
	}
	b, err := io.ReadAll(req.Body)
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
	short, duplicate := h.service.Storage.SetShort(j.URL)
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
}

func (h *Handler) GetUserURLAction(res http.ResponseWriter, req *http.Request) {
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
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	for _, short := range h.service.Storage.GetFullList() {
		if GetSignerCheck(short.Signer.Sign, cookie.Value) {
			obj := JSONObject{}
			obj.ShortURL = short.ShortURL
			obj.OriginalURL = short.LongURL
			list = append(list, obj)
		}

	}
	if len(list) == 0 {
		http.Error(res, "No content!", http.StatusNoContent)

		return
	}
	p, _ := json.Marshal(list)
	res.Write([]byte(p))
}

func (h *Handler) GetPingAction(res http.ResponseWriter, req *http.Request) {
	if status, err := ConnectionDBCheck(); status != http.StatusOK {
		http.Error(res, err, http.StatusInternalServerError)

		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}

func (h *Handler) GetBatchAction(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed for this route!", http.StatusMethodNotAllowed)

		return
	}
	if req.URL.Path != "/api/shorten/batch" {
		http.Error(res, "Wrong route!", http.StatusNotFound)

		return
	}
	defer req.Body.Close()
	if req.ContentLength < 1 {
		http.Error(res, "Empty body!", http.StatusBadRequest)
	}
	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Add("Accept", "application/json")
	list := []JSONBatcher{}
	if err := json.Unmarshal(b, &list); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	resultsObj := []JSONResultBatcher{}
	for i, obj := range list {
		short, _ := h.service.Storage.SetShort(obj.LongURL)
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
		resultBatcher.URLID = obj.URLID
		resultBatcher.ShortURL = short.ShortURL
		resultsObj = append(resultsObj, *resultBatcher)
	}
	res.WriteHeader(http.StatusCreated)
	if len(resultsObj) == 0 {
		http.Error(res, "No content!", http.StatusNoContent)

		return
	}
	p, _ := json.Marshal(resultsObj)
	res.Write([]byte(p))
}

func (h *Handler) RemoveBatchAction(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(res, "Only Delete requests are allowed for this route!", http.StatusMethodNotAllowed)

		return
	}
	if req.URL.Path != "/api/user/urls" {
		http.Error(res, "Wrong route!", http.StatusNotFound)

		return
	}
	defer req.Body.Close()
	if req.ContentLength < 1 {
		http.Error(res, "Empty body!", http.StatusBadRequest)
	}
	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	cookie, _ := req.Cookie("user_id")
	if cookie == nil {
		http.Error(res, "Failed to identify, no 'user_id' cookie set", http.StatusBadRequest)
	} else {
		var list []string
		var shorters []string
		if err := json.Unmarshal(b, &list); err != nil { // тут может быть ошибка если будет передаваться не в json
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
		inputCh := make(chan *service.Shorter)

		go RemoveWorkers(
			h,
			list,
			cookie.Value,
			inputCh,
			shorters,
		)

	}
	res.WriteHeader(http.StatusAccepted)
}

func RemoveWorkers(
	h *Handler,
	list []string,
	userId string,
	inputCh chan *service.Shorter,
	shorters []string) {

	go func() {
		for _, id := range list {
			inputCh <- h.service.Storage.GetShorter(id)
		}
		close(inputCh)
	}()
	workersCount := 5
	workerChs := make([]chan *service.Shorter, 0, workersCount)
	fanOutChs := fanOut(inputCh, workersCount)
	for _, fanOutCh := range fanOutChs {
		workerCh := make(chan *service.Shorter)
		newWorker(fanOutCh, workerCh, userId)
		workerChs = append(workerChs, workerCh)
	}
	for id := range fanIn(workerChs...) {
		shorters = append(shorters, id)
	}
	h.service.Storage.RemoveShorts(shorters)
}

func fanOut(inputCh chan *service.Shorter, n int) []chan *service.Shorter {
	chs := make([]chan *service.Shorter, 0, n)
	for i := 0; i < n; i++ {
		ch := make(chan *service.Shorter)
		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan *service.Shorter) {
			for _, ch := range chs {
				close(ch)
			}
		}(chs)

		for i := 0; ; i++ {
			if i == len(chs) {
				i = 0
			}

			list, ok := <-inputCh
			if !ok {
				return
			}
			ch := chs[i]
			ch <- list
		}
	}()

	return chs
}

func newWorker(input, out chan *service.Shorter, userID string) {
	go func() {
		for shorter := range input {
			if GetSignerCheck(shorter.Signer.Sign, userID) {
				out <- shorter
			}

		}

		close(out)
	}()
}

func fanIn(inputChs ...chan *service.Shorter) chan string {
	outCh := make(chan string)
	go func() {
		wg := &sync.WaitGroup{}
		for _, inputCh := range inputChs {
			wg.Add(1)
			go func(inputCh chan *service.Shorter) {
				defer wg.Done()
				for shorter := range inputCh {
					outCh <- shorter.ID
				}
			}(inputCh)
		}
		wg.Wait()
		close(outCh)
	}()

	return outCh
}
