package token

import (
	"crypto/rand"
	"crypto/rsa"
	"gofemart/internal/models"
	"strconv"
	"testing"
	"time"
)

func TestGenerateAndParse(t *testing.T) {
	pkey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &pkey.PublicKey

	tests := []struct {
		name      string
		generator *JWTGenerator
		user      *models.User
		wantErr   bool
	}{
		{
			name:      "valid",
			generator: NewJWTGenerator(pkey, pubKey, time.Hour),
			user:      &models.User{ID: 1},
			wantErr:   false,
		},
		{
			name:      "expired",
			generator: NewJWTGenerator(pkey, pubKey, -time.Hour),
			user:      &models.User{ID: 2},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := tt.generator.Generate(tt.user)
			if err != nil {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			token, err := tt.generator.Parse(tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				subject, errClaims := token.Claims.GetSubject()
				if errClaims != nil {
					t.Errorf("Expected Subject, got error %s", errClaims.Error())
					return
				}
				intSubject, errParse := strconv.ParseInt(subject, 10, 64)
				if errParse != nil {
					t.Errorf("Error while subject to int = %v, want %v. Error %s", subject, tt.user.ID, errParse.Error())
					return
				}
				if intSubject != tt.user.ID {
					t.Errorf("Parse().Subject = %v, want %v", subject, tt.user.ID)
				}
			}
		})
	}
}
