package config

import (
	"crypto/rsa"
	"time"
)

const (
	// DefaultServerURL Url сервера получателя метрик по умолчанию
	DefaultServerURL = "localhost:8080"
	// DefaultLogLevel Уровень логирования по умолчанию
	DefaultLogLevel = "info"
	// DefaultDatabaseDSN подключение к базе данных
	DefaultDatabaseDSN = "postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"
	// DefaultHashKey ключ шифрования по умолчанию
	DefaultHashKey = "dsbtyrew3!hgsdfvbytrd324"
	// DefaultAccrualSystemAddress адрес системы расчёта начислений по умолчанию
	DefaultAccrualSystemAddress = "http://localhost:8480"
	// DefaultPrivateKeyPath Путь к приватному ключу для JWT по умолчанию
	DefaultPrivateKeyPath = ""
	// DefaultPublicKeyPath Путь к публичному ключу для JWT по умолчанию
	DefaultPublicKeyPath = ""
	// DefaultTokenExpiration Время жизни токена авторизации по умолчанию
	DefaultTokenExpiration = 12 * time.Hour

	// DefaultPrivateKey Текстовое представление приватного ключа для JWT по умолчанию
	DefaultPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgF8sjeK3PDamGn0icENKSjwWpuWjmPrKEoVWXIO7os5iM1CZOn7h
qC4OgRKArfNC2BVa2zvVcrzRxRzFobyM6fblMbRzgE//5ct6YtpkWUgWEOZfXZ/X
FY5AUlngBKZtU2MS/CUX+PFXICUIVTDoCL6ngwNqTQj/dCin6E75Q1c/AgMBAAEC
gYAnYTsQHPs4LYB2WIKVBS80L7c8+3U4B9aj/zjmdQQHW1CaP9yZVWuOKwgzDLVt
GzJnm6Fs34PLJwzlO80RREkmEynnTYNVJejCLwuyT1oEGV6rFsql2HcIZ073NCxN
WakwL6Ay7QHH5S+hJDHCuxAx7kKoqiIXRcvbcwpRAnE5kQJBAKISE5uw1ejUocnE
ad7M2PVTz36ZS9d/3glpRQiQ2exeFRtcsq1J6O7G5OK62UMv3tcHyjf2suY0fPAA
jlPt/jcCQQCWVTylcs1Q319VRecJxSiCPjj97AA2VO1gcgzCWQ7mTp+N8QIegrD/
ZvvHqSLt79CexWrnOI6SvPuMf+8fwas5AkAlw76L7cW6bimQ4VKmFueLKs9TuZbB
jUsIuF3cpBwThsy2RoBf/rPnR7M33cAYdsQfKPKG3dZL6/kc15RSnEc7AkAgXTdS
MxXqjDw84nCr1Ms0xuqEF/Ovvrbf5Y3DpWKkyFZnO3SGVwJ96ZDY2hvP96oFFGFA
aBehlZfeFojHYG1ZAkEAnfwWAoPmvHxDaakOMsZg9PVVHIMhJ3Uck7lU5HKofHhq
rW4FGtaAhyoIZ2DQgctfe+PMcflOzkzkg9Cpqax7Cg==
-----END RSA PRIVATE KEY-----`
	// DefaultPublicKey Текстовое представление публичного ключа для JWT по умолчанию
	DefaultPublicKey = `-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgF8sjeK3PDamGn0icENKSjwWpuWj
mPrKEoVWXIO7os5iM1CZOn7hqC4OgRKArfNC2BVa2zvVcrzRxRzFobyM6fblMbRz
gE//5ct6YtpkWUgWEOZfXZ/XFY5AUlngBKZtU2MS/CUX+PFXICUIVTDoCL6ngwNq
TQj/dCin6E75Q1c/AgMBAAE=
-----END PUBLIC KEY-----`
)

type JWTKeys struct {
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

// CliConfig конфигурация сервера из командной строки
type CliConfig struct {
	Address              string        `env:"RUN_ADDRESS"`            // адрес сервера
	LogLevel             string        `env:"LOG_LEVEL"`              // Уровень логирования
	DatabaseDSN          string        `env:"DATABASE_URI"`           // подключение к базе данных
	AccrualSystemAddress string        `env:"ACCRUAL_SYSTEM_ADDRESS"` // адрес системы расчёта начислений
	HashKey              string        `env:"KEY"`                    // Ключ для шифрования
	PrivateKeyPath       string        `env:"PKEYP"`                  // Путь к приватному ключу для JWT
	PublicKeyPath        string        `env:"PUKEYP"`                 // Путь к публичному ключу для JWT
	PrivateKey           string        `env:"PKEY"`                   // Приватный ключ для JWT
	PublicKey            string        `env:"PUKEY"`                  // Публичный ключ для JWT
	JWTKeys              *JWTKeys      `env:"-"`                      // Ключи для JWT
	TokenExpiration      time.Duration `env:"TOKEN_EXPIRATION"`       // Время жизни токена авторизации
}

// Params конфигурация приложения
var Params *CliConfig

// InitializeDefaultConfig инициализация конфигурации приложения
func InitializeDefaultConfig() *CliConfig {
	return &CliConfig{
		Address:              DefaultServerURL,
		LogLevel:             DefaultLogLevel,
		DatabaseDSN:          DefaultDatabaseDSN,
		AccrualSystemAddress: DefaultAccrualSystemAddress,
		HashKey:              DefaultHashKey,
		PrivateKeyPath:       DefaultPrivateKeyPath,
		PublicKeyPath:        DefaultPublicKeyPath,
		TokenExpiration:      DefaultTokenExpiration,
		PrivateKey:           DefaultPrivateKey,
		PublicKey:            DefaultPublicKey,
	}
}
