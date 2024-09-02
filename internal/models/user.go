package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	config "gofemart/internal/configuration"
)

type User struct {
	Id           int64  `db:"id"`
	Login        string `db:"login"`
	Password     string `db:"-"`
	PasswordHash string `db:"password_hash"`
}

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
