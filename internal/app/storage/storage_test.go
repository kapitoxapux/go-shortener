package storage

import (
	"net/http"
	"os"
	"testing"
	"time"

	"myapp/internal/app/config"
	"myapp/internal/app/handler"
	"myapp/internal/app/repository"
	"myapp/internal/app/service"
)

var repo repository.Repository

func GetDB() service.Storage {

	if status, _ := handler.ConnectionDBCheck(); status == http.StatusOK {
		return NewDB()
	}

	if pathStorage := config.GetConfigPath(); pathStorage != "" {
		return NewFileDB()
	}

	return NewInMemDB()
}

func Test_getFullUrl(t *testing.T) {
	if status, _ := handler.ConnectionDBCheck(); status == http.StatusOK {
		repo = repository.NewRepository(config.GetStorageDB())
	} else {
		repo = nil
	}

	db := GetDB()
	s := service.NewService(db)

	dataForCookie := time.Now().String()
	cookie := handler.SetCookieToken(dataForCookie)
	s2, _ := s.Storage.SetShort(os.Getenv("BASE_URL")+"/some_text_to_test_2", cookie)
	tests := []struct {
		name    string
		link    string
		shorter *service.Shorter
		equel   bool
	}{
		{
			name:    "unique check",
			link:    os.Getenv("BASE_URL") + "/some_text_to_test_2",
			shorter: s2,
			equel:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.Storage.GetFullURL(tt.link); (got != tt.shorter.ShortURL) == tt.equel {
				t.Errorf("getFullUrl() = %v, want %v", got, tt.equel)
			}
		})
	}
}

func Test_getShort(t *testing.T) {
	if status, _ := handler.ConnectionDBCheck(); status == http.StatusOK {
		repo = repository.NewRepository(config.GetStorageDB())
	} else {
		repo = nil
	}

	db := GetDB()
	s := service.NewService(db)

	dataForCookie := time.Now().String()
	cookie := handler.SetCookieToken(dataForCookie)
	s1, _ := s.Storage.SetShort(os.Getenv("BASE_URL")+"/some_text_to_test_1", cookie)
	tests := []struct {
		name    string
		link    string
		shorter *service.Shorter
		equel   bool
	}{
		{
			name:    "unique check",
			link:    os.Getenv("BASE_URL") + "/some_text_to_test_1",
			shorter: s1,
			equel:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.Storage.GetShort(tt.link); (got != tt.shorter.ShortURL) == tt.equel {
				t.Errorf("getShort() = %v, want %v", got, tt.equel)
			}
		})
	}
}

func Test_setShort(t *testing.T) {

	if status, _ := handler.ConnectionDBCheck(); status == http.StatusOK {
		repo = repository.NewRepository(config.GetStorageDB())
	} else {
		repo = nil
	}

	db := GetDB()
	s := service.NewService(db)

	dataForCookie := time.Now().String()
	cookie := handler.SetCookieToken(dataForCookie)

	testNegative, _ := s.Storage.SetShort(os.Getenv("BASE_URL")+"/some_text_to_test_1", cookie)
	testPositive, _ := s.Storage.SetShort(os.Getenv("BASE_URL")+"/some_text_to_test_2", cookie)

	tests := []struct {
		name    string
		link    string
		want    *service.Shorter
		wantErr bool
	}{
		{
			name:    "new Shorter",
			link:    os.Getenv("BASE_URL") + "/some_text_to_test_2",
			want:    testPositive,
			wantErr: false,
		},
		{
			name:    "catch error",
			link:    os.Getenv("BASE_URL") + "/some_text_to_test_1",
			want:    testNegative,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testPositive; (got.ID != tt.want.ID) != tt.wantErr {
				t.Errorf("setShort() = %v, want %v", got, tt.want.ID)
			}
		})
	}
}
