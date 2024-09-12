package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	config "gofemart/internal/configuration"
)

type User struct {
	ID           int64  `db:"id"`
	Login        string `db:"login"`
	Password     string `db:"-"`
	PasswordHash string `db:"password_hash"`
}

// GeneratePasswordHash создаём хэш пароля пользователя
func (u *User) GeneratePasswordHash() error {
	if u.Password == "" {
		return errors.New("password key is empty")
	}
	if config.Params.HashKey == "" {
		return errors.New("hash key is empty")
	}
	harsher := hmac.New(sha256.New, []byte(config.Params.HashKey))
	harsher.Write([]byte(u.Password))
	u.PasswordHash = hex.EncodeToString(harsher.Sum(nil))

	return nil
}

// CheckPasswordHash проверяем совпадение паролей
func (u *User) CheckPasswordHash(passwordHash string) (bool, error) {
	decodedHash, err := hex.DecodeString(passwordHash)
	if err != nil {
		return false, err
	}
	decodedUserHash, err := hex.DecodeString(u.PasswordHash)
	if err != nil {
		return false, err
	}

	return hmac.Equal(decodedHash, decodedUserHash), nil
}
