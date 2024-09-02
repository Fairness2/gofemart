package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gofemart/internal/models"
	"strconv"
	"time"
)

// JWTGenerator класс использования jwt токенов
type JWTGenerator struct {
	pkey       *rsa.PrivateKey
	pubKey     *rsa.PublicKey
	expiration time.Duration
	method     jwt.SigningMethod
	issuer     string
}

// NewJWTGenerator создание нового генератора
func NewJWTGenerator(pkey *rsa.PrivateKey, pubKey *rsa.PublicKey, expiration time.Duration) *JWTGenerator {
	return &JWTGenerator{
		pkey:       pkey,
		pubKey:     pubKey,
		expiration: expiration,
		method:     jwt.SigningMethodRS256,
		issuer:     "gofemart",
	}
}

// Generate создание нового токена для пользователя
func (g *JWTGenerator) Generate(user *models.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    g.issuer,
		Subject:   strconv.FormatInt(user.Id, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.expiration)),
	}

	token := jwt.NewWithClaims(g.method, claims)
	return token.SignedString(g.pkey)
}

// Parse парсим полученный токен
func (g *JWTGenerator) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return g.pubKey, nil
		},
		jwt.WithIssuer(g.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{"RS256"}),
	)
}
