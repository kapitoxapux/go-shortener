package storage

import (
	"net/http"
	"os"
	"testing"

	"myapp/internal/app/config"
	"myapp/internal/app/handler"
	"myapp/internal/app/repository"
	"myapp/internal/app/service"
)

type Service struct {
	Storage service.Storage
}

func NewService() *Service {
	return &Service{}
}

var Srv = NewService()

var repo repository.Repository

func Test_getFullUrl(t *testing.T) {
	if status, _ := handler.ConnectionDBCheck(); status == http.StatusOK {
		repo = repository.NewRepository(config.GetStorageDB())
	} else {
		repo = nil
	}

	s2, _ := Srv.Storage.SetShort(os.Getenv("BASE_URL") + "/some_text_to_test_2")
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
			if got := Srv.Storage.GetFullURL(tt.link); (got != tt.shorter.ShortURL) == tt.equel {
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

	s1, _ := Srv.Storage.SetShort(os.Getenv("BASE_URL") + "/some_text_to_test_1")
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
			if got := Srv.Storage.GetShort(tt.link); (got != tt.shorter.ShortURL) == tt.equel {
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

	testNegative, _ := Srv.Storage.SetShort(os.Getenv("BASE_URL") + "/some_text_to_test_1")
	testPositive, _ := Srv.Storage.SetShort(os.Getenv("BASE_URL") + "/some_text_to_test_2")

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
