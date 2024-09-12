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
	DefaultHashKey = ""
	// DefaultAccrualSystemAddress адрес системы расчёта начислений по умолчанию
	DefaultAccrualSystemAddress = "http://localhost:8480"
	// DefaultPrivateKeyPath Путь к приватному ключу для JWT по умолчанию
	DefaultPrivateKeyPath = ""
	// DefaultPublicKeyPath Путь к публичному ключу для JWT по умолчанию
	DefaultPublicKeyPath = ""
	// DefaultTokenExpiration Время жизни токена авторизации по умолчанию
	DefaultTokenExpiration = 12 * time.Hour
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
	PrivateKeyPath       string        `env:"PKEY"`                   // Путь к приватному ключу для JWT
	PublicKeyPath        string        `env:"PUKEY"`                  // Путь к публичному ключу для JWT
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
	}
}
