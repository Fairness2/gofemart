package token

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

// ParseKeys получаем ключи для JWT токенов
func ParseKeys(privateKeyPath string, publicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if privateKeyPath == "" {
		return nil, nil, errors.New("no private key path specified")
	}
	if publicKeyPath == "" {
		return nil, nil, errors.New("no public key path specified")
	}

	pkeyBody, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}
	pkey, err := jwt.ParseRSAPrivateKeyFromPEM(pkeyBody)
	if err != nil {
		return nil, nil, err
	}

	pubKeyBody, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBody)
	if err != nil {
		return nil, nil, err
	}

	return pkey, pubKey, nil
}
