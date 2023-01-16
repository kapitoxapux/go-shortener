package storage

import (
	"os"
	"testing"
)

func Test_getFullUrl(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		shorter *Shorter
		equel   bool
	}{
		{
			name:    "unique check",
			link:    "http://" + os.Getenv("SERVER_ADDRESS") + os.Getenv("BASE_URL") + "/some_text_to_test_2",
			shorter: SetShort("http://" + os.Getenv("SERVER_ADDRESS") + os.Getenv("BASE_URL") + "/some_text_to_test_2"),
			equel:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFullURL(tt.link); (got != tt.shorter.ShortURL) == tt.equel {
				t.Errorf("getFullUrl() = %v, want %v", got, tt.equel)
			}
		})
	}
}

func Test_getShort(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		shorter *Shorter
		equel   bool
	}{
		{
			name:    "unique check",
			link:    "http://" + os.Getenv("SERVER_ADDRESS") + os.Getenv("BASE_URL") + "/some_text_to_test_1",
			shorter: SetShort("http://" + os.Getenv("SERVER_ADDRESS") + os.Getenv("BASE_URL") + "/some_text_to_test_1"),
			equel:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetShort(tt.link); (got != tt.shorter.ShortURL) == tt.equel {
				t.Errorf("getShort() = %v, want %v", got, tt.equel)
			}
		})
	}
}

func Test_setShort(t *testing.T) {

	testNegative := SetShort("http://" + os.Getenv("SERVER_ADDRESS") + os.Getenv("BASE_URL") + "/some_text_to_test_1")
	testPositive := SetShort("http://" + os.Getenv("SERVER_ADDRESS") + os.Getenv("BASE_URL") + "/some_text_to_test_2")

	tests := []struct {
		name    string
		link    string
		want    *Shorter
		wantErr bool
	}{
		{
			name:    "new Shorter",
			link:    "http://" + os.Getenv("SERVER_ADDRESS") + os.Getenv("BASE_URL") + "/some_text_to_test_2",
			want:    testPositive,
			wantErr: false,
		},
		{
			name:    "catch error",
			link:    "http://" + os.Getenv("SERVER_ADDRESS") + os.Getenv("BASE_URL") + "/some_text_to_test_1",
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
