package config

import (
	"crypto/rsa"
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

// NewConfig инициализирует новую консольную конфигурацию, обрабатывает аргументы командной строки
func NewConfig() (*CliConfig, error) {
	// Регистрируем новое хранилище
	cnf := NewDefaultConfig()
	// Заполняем конфигурацию из параметров командной строки
	if err := parseFromCli(cnf); err != nil {
		return nil, err
	}
	// Заполняем конфигурацию из окружения
	if err := parseFromEnv(cnf); err != nil {
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
	return nil
}

// parseFromCli заполняем конфигурацию из параметров командной строки
func parseFromCli(cnf *CliConfig) error {
	// Регистрируем флаги конфигурации
	flag.StringVar(&cnf.Address, "a", DefaultServerURL, "address and port to run server")
	flag.StringVar(&cnf.LogLevel, "ll", DefaultLogLevel, "level of logging")
	flag.StringVar(&cnf.DatabaseDSN, "d", DefaultDatabaseDSN, "database connection")
	flag.StringVar(&cnf.AccrualSystemAddress, "к", DefaultAccrualSystemAddress, "accrual system address")
	flag.StringVar(&cnf.HashKey, "k", DefaultHashKey, "encrypted key")
	flag.StringVar(&cnf.PrivateKeyPath, "pkp", DefaultPrivateKeyPath, "path to private key")
	flag.StringVar(&cnf.PublicKeyPath, "pukp", DefaultPublicKeyPath, "path to public key")
	flag.StringVar(&cnf.PrivateKey, "pk", DefaultPrivateKey, "private key")
	flag.StringVar(&cnf.PublicKey, "puk", DefaultPublicKey, "public key")

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
