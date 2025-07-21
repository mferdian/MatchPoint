package helpers

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", "mySecret123", false},
		{"Empty password", "", false}, // bcrypt tetap bisa hash empty string
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("HashPassword() got empty hash")
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	plainPassword := "mySecret123"
	hash, err := HashPassword(plainPassword)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	tests := []struct {
		name          string
		hashPassword  string
		plainPassword []byte
		want          bool
		wantErr       bool
	}{
		{"Correct password", hash, []byte(plainPassword), true, false},
		{"Wrong password", hash, []byte("wrongPassword"), false, true},
		{"Empty password", hash, []byte(""), false, true},
		{"Invalid hash format", "invalidHash", []byte(plainPassword), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckPassword(tt.hashPassword, tt.plainPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
