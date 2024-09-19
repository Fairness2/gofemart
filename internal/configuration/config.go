package config

import (
	"crypto/rsa"
	_ "embed"
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
	// DefaultAccrualSenderPause пауза в запросах к сервису начислений, если он ответил ответом, что слишком много запросов
	DefaultAccrualSenderPause = time.Minute
	// DefaultQueueSize количество заказов, которые одновременно могут находиться в очереди на проверке, если очередь заполнена, то они будут отложены
	DefaultQueueSize = 1000
	// DefaultWorkerCount количество обработчиков заказов
	DefaultWorkerCount = 10
	// DefaultDBCheckDuration период в который проверяется база данных на необработанные заказы
	DefaultDBCheckDuration = 5 * time.Second
)

// DefaultPrivateKey Текстовое представление приватного ключа для JWT по умолчанию
//
//go:embed keys/private.pem
var DefaultPrivateKey string

// DefaultPublicKey Текстовое представление публичного ключа для JWT по умолчанию
//
//go:embed keys/public.pem
var DefaultPublicKey string

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
	AccrualSenderPause   time.Duration `env:"ACCRUAL_SENDER_PAUSE"`   // пауза в запросах к сервису начислений, если он ответил ответом, что слишком много запросов
	QueueSize            int           `env:"QUEUE_SIZE"`             // количество заказов, которые одновременно могут находиться в очереди на проверке, если очередь заполнена, то они будут отложены
	WorkerCount          int           `env:"WORKER_COUNT"`           // количество обработчиков заказов
	DBCheckDuration      time.Duration `env:"DB_CHECK_DURATION"`      // период в который проверяется база данных на необработанные заказы
}

// NewDefaultConfig инициализация конфигурации приложения
func NewDefaultConfig() *CliConfig {
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
		AccrualSenderPause:   DefaultAccrualSenderPause,
		QueueSize:            DefaultQueueSize,
		WorkerCount:          DefaultWorkerCount,
		DBCheckDuration:      DefaultDBCheckDuration,
	}
}
