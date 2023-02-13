package storage

// "bufio"

// "math/big"
// "math/rand"
// "os"
// "time"

// _ "github.com/jackc/pgx/v5/stdlib"
// "myapp/internal/app/models"
// "myapp/internal/app/service"

// var short string = ""

// var paths = map[string]*Shorter{}

// type Signer struct {
// 	SignID uint32 `json:"signID"`
// 	Sign   []byte `json:"sign"`
// }

// type Shorter struct {
// 	ID       string
// 	LongURL  string `json:"longURL"`
// 	ShortURL string `json:"shortURL"`
// 	BaseURL  string `json:"baseURL"`
// 	Signer
// }

// func NewShorter() Shorter {
// 	shorter := Shorter{}
// 	shorter.ID = ""
// 	shorter.LongURL = ""
// 	shorter.ShortURL = ""
// 	shorter.BaseURL = config.GetConfigBase() + "/"
// 	shorter.Signer.SignID = 0
// 	shorter.Signer.Sign = nil

// 	return shorter
// }

// func ShorterSignerSet(short string) Signer {
// 	data, _ := hex.DecodeString(short)
// 	h := hmac.New(sha256.New, config.Secretkey)
// 	h.Write(data)
// 	sign := h.Sum(nil)
// 	id := binary.BigEndian.Uint32(sign[:4])

// 	return Signer{id, sign}
// }

// func ConnectionDBCheck() (int, string) {
// 	db, err := sql.Open("pgx", config.GetStorageDB())
// 	if err != nil {

// 		return 500, err.Error()
// 	}

// 	// close database
// 	defer db.Close()

// 	// check db
// 	err = db.Ping()
// 	if err != nil {

// 		return 500, err.Error()
// 	}

// 	return 200, ""
// }

// func (s *service.Storage) SetShort(link string) (*Shorter, bool) {
// 	shorter := NewShorter()
// 	duplicate := false

// 	shorter = s.storage.ShowShortenerByLong(link)

// 	return &shorter, duplicate
// }

// func GetShort(repo repository.Repository, id string) string {
// 	shortURL := ""
// 	if repo != nil {
// 		if result, err := repo.ShowShortener(id); err != nil {
// 			log.Fatal("Короткая ссылка не найдена, произошла ошибка: %w", err)
// 		} else {
// 			shortURL = result.ShortURL
// 		}
// 	} else {
// 		pathStorage := config.GetConfigPath()
// 		if pathStorage == "" {
// 			if paths[id] != nil {

// 				return paths[id].ShortURL
// 			}

// 		} else {
// 			reader, _ := NewReader(pathStorage)
// 			defer reader.Close()

// 			shorter := NewShorter()
// 			for reader.scanner.Scan() {
// 				data := reader.scanner.Bytes()

// 				_ = json.Unmarshal(data, &shorter)
// 				if id == shorter.ID {
// 					return shorter.ShortURL
// 				}

// 			}

// 		}
// 	}

// 	return shortURL
// }

// func GetFullURL(repo repository.Repository, id string) string {
// 	longURL := ""
// 	if repo != nil {
// 		if result, err := repo.ShowShortener(id); err != nil {
// 			log.Fatal("Полная ссылка не найдена, произошла ошибка: %w", err)
// 		} else {
// 			longURL = result.LongURL
// 		}
// 	} else {
// 		pathStorage := config.GetConfigPath()
// 		if pathStorage == "" {
// 			if paths[id] != nil {
// 				return paths[id].LongURL
// 			}

// 		} else {
// 			reader, _ := NewReader(pathStorage)
// 			defer reader.Close()

// 			shorter := NewShorter()
// 			for reader.scanner.Scan() {
// 				data := reader.scanner.Bytes()

// 				_ = json.Unmarshal(data, &shorter)
// 				if id == shorter.ID {
// 					return shorter.LongURL
// 				}

// 			}

// 		}
// 	}

// 	return longURL
// }

// func GetFullList(repo repository.Repository) map[string]*Shorter {
// 	if repo != nil {
// 		if results, err := repo.ShowShorteners(); err != nil {
// 			log.Fatal("Произошла ошибка получения списка: %w", err)
// 		} else {
// 			for _, model := range results {
// 				shorter := NewShorter()
// 				shorter.ID = model.ID
// 				shorter.ShortURL = model.ShortURL
// 				shorter.LongURL = model.LongURL
// 				shorter.Sign = model.Sign
// 				shorter.SignID = model.SignID

// 				paths[model.ID] = &shorter
// 			}
// 		}
// 	} else {
// 		if pathStorage := config.GetConfigPath(); pathStorage != "" {
// 			reader, _ := NewReader(pathStorage)
// 			defer reader.Close()

// 			for reader.scanner.Scan() {
// 				data := reader.scanner.Bytes()

// 				shorter := NewShorter()
// 				_ = json.Unmarshal(data, &shorter)
// 				paths[shorter.ID] = &shorter
// 			}

// 		}
// 	}

// 	return paths
// }
