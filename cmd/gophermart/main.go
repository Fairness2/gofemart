package main

import (
	"gofemart/internal/application"
	config "gofemart/internal/configuration"
	"gofemart/internal/logger"
	"log"
)

// @title GoFemart API
// @version 1.0
// @description Система баланса поощрений

// @host localhost:8080
// @BasePath /
func main() {
	log.Println("Start program")
	// Устанавливаем настройки
	cnf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	_, err = logger.NewGlobal(cnf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	// Показываем конфигурацию сервера
	logger.Log.Infow("Running server with configuration",
		"address", cnf.Address,
		"logLevel", cnf.LogLevel,
		"databaseDSN", cnf.DatabaseDSN,
		"accrualSystemAddress", cnf.AccrualSystemAddress,
	)

	// стартуем приложение
	if err = application.New(cnf); err != nil {
		logger.Log.Error(err)
	}

	logger.Log.Info("End program")
}
