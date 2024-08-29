package config

const (
	// DefaultServerURL Url сервера получателя метрик по умолчанию
	DefaultServerURL = "localhost:8080"
	// DefaultLogLevel Уровень логирования по умолчанию
	DefaultLogLevel = "info"
	// DefaultDatabaseDSN подключение к базе данных
	DefaultDatabaseDSN = "" //"postgresql://postgres:example@127.0.0.1:5432/gmetrics"
	// DefaultHashKey ключ шифрования по умолчанию
	DefaultHashKey = ""
	// DefaultAccrualSystemAddress адрес системы расчёта начислений по умолчанию
	DefaultAccrualSystemAddress = "localhost:8480"
)

// CliConfig конфигурация сервера из командной строки
type CliConfig struct {
	Address              string `env:"RUN_ADDRESS"`            // адрес сервера
	LogLevel             string `env:"LOG_LEVEL"`              // Уровень логирования
	DatabaseDSN          string `env:"DATABASE_URI"`           // подключение к базе данных
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"` // адрес системы расчёта начислений
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
	}
}
