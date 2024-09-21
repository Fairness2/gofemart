package config

import (
	"crypto/rsa"
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

// NewConfig инициализирует новую консольную конфигурацию, обрабатывает аргументы командной строки
func NewConfig() (*CliConfig, error) {

	/*// Регистрируем новое хранилище
	cnf := NewDefaultConfig()
	// Заполняем конфигурацию из параметров командной строки
	if err := parseFromCli(cnf); err != nil {
		return nil, err
	}
	// Заполняем конфигурацию из окружения
	if err := parseFromEnv(cnf); err != nil {
		return nil, err
	}*/
	cnf := &CliConfig{}
	if err := parseFromViper(cnf); err != nil {
		return nil, err
	}
	// Парсим ключи для JWT токена
	if err := parseKeys(cnf); err != nil {
		return nil, err
	}

	return cnf, nil
}

// parseFromEnv заполняем конфигурацию переменных из окружения
func parseFromEnv(params *CliConfig) error {
	cnf := CliConfig{}
	err := env.Parse(&cnf)
	// Если ошибка, то считаем, что вывести конфигурацию из окружения не удалось
	if err != nil {
		return err
	}
	if cnf.Address != "" {
		params.Address = cnf.Address
	}
	if cnf.LogLevel != "" {
		params.LogLevel = cnf.LogLevel
	}
	if cnf.DatabaseDSN != "" {
		params.DatabaseDSN = cnf.DatabaseDSN
	}
	if cnf.AccrualSystemAddress != "" {
		params.AccrualSystemAddress = cnf.AccrualSystemAddress
	}
	if cnf.HashKey != "" {
		params.HashKey = cnf.HashKey
	}
	if cnf.PrivateKeyPath != "" {
		params.PrivateKeyPath = cnf.PrivateKeyPath
	}
	if cnf.PublicKeyPath != "" {
		params.PublicKeyPath = cnf.PublicKeyPath
	}
	if cnf.PrivateKey != "" {
		params.PrivateKey = cnf.PrivateKey
	}
	if cnf.PublicKey != "" {
		params.PublicKey = cnf.PublicKey
	}
	if cnf.QueueSize > 0 {
		params.QueueSize = cnf.QueueSize
	}
	if cnf.WorkerCount > 0 {
		params.WorkerCount = cnf.WorkerCount
	}
	if cnf.TokenExpiration > 0 {
		params.TokenExpiration = cnf.TokenExpiration
	}
	if cnf.AccrualSenderPause > 0 {
		params.AccrualSenderPause = cnf.AccrualSenderPause
	}
	if cnf.DBCheckDuration > 0 {
		params.DBCheckDuration = cnf.DBCheckDuration
	}
	return nil
}

// parseFromCli заполняем конфигурацию из параметров командной строки
func parseFromCli(cnf *CliConfig) error {
	// Регистрируем флаги конфигурации
	flag.StringVar(&cnf.Address, "a", DefaultServerURL, "address and port to run server")
	flag.StringVar(&cnf.LogLevel, "ll", DefaultLogLevel, "level of logging")
	flag.StringVar(&cnf.DatabaseDSN, "d", DefaultDatabaseDSN, "database connection")
	flag.StringVar(&cnf.AccrualSystemAddress, "к", DefaultAccrualSystemAddress, "accrual system address")
	flag.StringVar(&cnf.HashKey, "hk", DefaultHashKey, "encrypted key")
	flag.StringVar(&cnf.PrivateKeyPath, "pkp", DefaultPrivateKeyPath, "path to private key")
	flag.StringVar(&cnf.PublicKeyPath, "pukp", DefaultPublicKeyPath, "path to public key")
	flag.StringVar(&cnf.PrivateKey, "pk", DefaultPrivateKey, "private key")
	flag.StringVar(&cnf.PublicKey, "puk", DefaultPublicKey, "public key")
	flag.IntVar(&cnf.QueueSize, "qs", DefaultQueueSize, "size of queue")
	flag.IntVar(&cnf.WorkerCount, "wc", DefaultWorkerCount, "count of workers")
	flag.DurationVar(&cnf.TokenExpiration, "te", DefaultTokenExpiration, "token expiration time")
	flag.DurationVar(&cnf.AccrualSenderPause, "aps", DefaultAccrualSenderPause, "pause between sending accruals")
	flag.DurationVar(&cnf.DBCheckDuration, "dbs", DefaultDBCheckDuration, "duration between BD checks")

	// Парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse() // Сейчас будет выход из приложения, поэтому код ниже не будет исполнен, но может пригодиться в будущем, если поменять флаг выхода или будет несколько сетов
	if !flag.Parsed() {
		return errors.New("error while parse flags")
	}
	return nil
}

// parseKeys парсим ключи для JWT токена
func parseKeys(cnf *CliConfig) error {
	pkey, pubKey, err := parseKeysFromFile(cnf.PrivateKeyPath, cnf.PublicKeyPath, cnf.PrivateKey, cnf.PublicKey)
	if err != nil {
		return err
	}
	cnf.JWTKeys = &JWTKeys{
		Public:  pubKey,
		Private: pkey,
	}
	return nil
}

// parseKeysFromFile получаем ключи для JWT токенов
func parseKeysFromFile(privateKeyPath string, publicKeyPath string, privateKey string, publicKey string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if privateKeyPath == "" && privateKey == "" {
		return nil, nil, errors.New("no private key path specified")
	}
	if publicKeyPath == "" && publicKey == "" {
		return nil, nil, errors.New("no public key path specified")
	}

	var pkeyBody []byte
	var pubKeyBody []byte
	var err error
	if privateKeyPath != "" {
		pkeyBody, err = os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, nil, err
		}
	} else {
		pkeyBody = []byte(privateKey)
	}

	pkey, err := jwt.ParseRSAPrivateKeyFromPEM(pkeyBody)
	if err != nil {
		return nil, nil, err
	}

	if publicKeyPath != "" {
		pubKeyBody, err = os.ReadFile(publicKeyPath)
		if err != nil {
			return nil, nil, err
		}
	} else {
		pubKeyBody = []byte(publicKey)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBody)
	if err != nil {
		return nil, nil, err
	}

	return pkey, pubKey, nil
}

// parseFromViper анализирует конфигурацию из переменных среды и аргументов командной строки с помощью Viper.
func parseFromViper(cnf *CliConfig) error {
	if err := bindEnv(); err != nil {
		return err
	}
	if err := bindArg(); err != nil {
		return err
	}

	return viper.Unmarshal(cnf)
}

// bindEnv привязывает переменные среды к ключам конфигурации Viper, гарантируя, что каждая привязка проверяется на наличие ошибок.
func bindEnv() error {
	if err := viper.BindEnv("Address", "RUN_ADDRESS"); err != nil {
		return err
	}
	if err := viper.BindEnv("LogLevel", "LOG_LEVEL"); err != nil {
		return err
	}
	if err := viper.BindEnv("DatabaseDSN", "DATABASE_URI"); err != nil {
		return err
	}
	if err := viper.BindEnv("AccrualSystemAddress", "ACCRUAL_SYSTEM_ADDRESS"); err != nil {
		return err
	}
	if err := viper.BindEnv("HashKey", "KEY"); err != nil {
		return err
	}
	if err := viper.BindEnv("PrivateKeyPath", "PKEYP"); err != nil {
		return err
	}
	if err := viper.BindEnv("PublicKeyPath", "PUKEYP"); err != nil {
		return err
	}
	if err := viper.BindEnv("PrivateKey", "PKEY"); err != nil {
		return err
	}
	if err := viper.BindEnv("PublicKey", "PUKEY"); err != nil {
		return err
	}
	if err := viper.BindEnv("TokenExpiration", "TOKEN_EXPIRATION"); err != nil {
		return err
	}
	if err := viper.BindEnv("AccrualSenderPause", "ACCRUAL_SENDER_PAUSE"); err != nil {
		return err
	}
	if err := viper.BindEnv("QueueSize", "QUEUE_SIZE"); err != nil {
		return err
	}
	if err := viper.BindEnv("WorkerCount", "WORKER_COUNT"); err != nil {
		return err
	}
	if err := viper.BindEnv("DBCheckDuration", "DB_CHECK_DURATION"); err != nil {
		return err
	}
	return nil
}

// bindArg привязывает аргументы командной строки к ключам конфигурации и устанавливает значения по умолчанию с помощью библиотек pflag и viper.
func bindArg() error {
	pflag.StringP("Address", "a", DefaultServerURL, "address and port to run server")
	pflag.StringP("LogLevel", "l", DefaultLogLevel, "level of logging")
	pflag.StringP("DatabaseDSN", "d", DefaultDatabaseDSN, "database connection")
	pflag.StringP("AccrualSystemAddress", "k", DefaultAccrualSystemAddress, "accrual system address")
	pflag.StringP("HashKey", "h", DefaultHashKey, "encrypted key")
	pflag.StringP("PrivateKeyPath", "p", DefaultPrivateKeyPath, "path to private key")
	pflag.StringP("PublicKeyPath", "u", DefaultPublicKeyPath, "path to public key")
	pflag.StringP("PrivateKey", "r", DefaultPrivateKey, "private key")
	pflag.StringP("PublicKey", "b", DefaultPublicKey, "public key")
	pflag.IntP("QueueSize", "q", DefaultQueueSize, "size of queue")
	pflag.IntP("WorkerCount", "w", DefaultWorkerCount, "count of workers")
	pflag.DurationP("TokenExpiration", "t", DefaultTokenExpiration, "token expiration time")
	pflag.DurationP("AccrualSenderPause", "s", DefaultAccrualSenderPause, "pause between sending accruals")
	pflag.DurationP("DBCheckDuration", "c", DefaultDBCheckDuration, "duration between BD checks")
	pflag.Parse()
	return viper.BindPFlags(pflag.CommandLine)
}
