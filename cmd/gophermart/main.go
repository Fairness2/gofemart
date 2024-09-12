package main

import (
	"context"
	"errors"
	"go.uber.org/zap"
	config "gofemart/internal/configuration"
	database "gofemart/internal/databse"
	"gofemart/internal/databse/migrations"
	"gofemart/internal/logger"
	"gofemart/internal/ordercheck"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Устанавливаем настройки
	cnf, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}
	config.Params = cnf

	_, err = initLogger()
	if err != nil {
		log.Fatal(err)
	}

	// Показываем конфигурацию сервера
	logger.Log.Infow("Running server with configuration",
		"address", config.Params.Address,
		"logLevel", config.Params.LogLevel,
		"databaseDSN", config.Params.DatabaseDSN,
		"accrualSystemAddress", config.Params.AccrualSystemAddress,
	)

	// стартуем приложение
	if err = runApplication(); err != nil {
		logger.Log.Error(err)
	}

	logger.Log.Info("End program")
}

// runApplication производим старт приложения
func runApplication() error {
	ctx, cancel := context.WithCancel(context.Background()) // Контекст для правильной остановки синхронизации
	defer func() {
		logger.Log.Info("Cancel context")
		cancel()
	}()

	// Вызываем функцию закрытия базы данных
	defer closeDB()
	// Инициализируем базу данных
	err := initDB()
	if err != nil {
		return err
	}

	ordercheck.CheckPool = ordercheck.NewPool(ctx, 1000, 10, time.Minute, config.Params.AccrualSystemAddress) // TODO конфигурация
	defer ordercheck.CheckPool.Close()

	wg := new(errgroup.Group)
	server := initServer()
	// Запускаем сервер
	wg.Go(func() error {
		sErr := server.ListenAndServe()
		if sErr != nil && !errors.Is(sErr, http.ErrServerClosed) {
			return sErr
		}
		return nil
	})
	// Регистрируем прослушиватель для закрытия записи в файл и завершения сервера
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logger.Log.Info("Stopping server")
	cancel()
	if err = stopServer(server, ctx); err != nil { // Запускаем сервер
		return err
	}

	// Ожидаем завершения всех горутин перед завершением программы
	if err = wg.Wait(); err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("End Server")
	return nil
}

// initLogger инициализируем логер
func initLogger() (*zap.SugaredLogger, error) {
	lgr, err := logger.New(config.Params.LogLevel)
	if err != nil {
		return nil, err
	}
	logger.Log = lgr

	return lgr, nil
}

// initDB инициализация подключения к бд
func initDB() error {
	// Создание пула подключений к базе данных для приложения
	var err error
	database.DB, err = database.NewPgDB(config.Params.DatabaseDSN)
	if err != nil {
		return err
	}

	if config.Params.DatabaseDSN != "" {
		logger.Log.Info("Migrate migrations")
		// Применим миграции
		migrator, err := migrations.New()
		if err != nil {
			return err
		}
		if err = migrator.Migrate(database.DB); err != nil {
			return err
		}
	}
	// создаём пул для SQLx на основе полученного пула стандартного SQL
	database.DBx = database.NewPgDBx(database.DB)

	return nil
}

// closeDB закрытие базы данных
func closeDB() {
	logger.Log.Info("Closing database connection for defer")
	if database.DB != nil {
		err := database.DB.Close()
		if err != nil {
			logger.Log.Error(err)
		}
	}
}

// initServer Создаём сервер приложения
func initServer() *http.Server {
	logger.Log.Infof("Running server on %s", config.Params.Address)
	server := http.Server{
		Addr:    config.Params.Address,
		Handler: getRouter(),
	}

	return &server
}

// stopServer закрытие сервера
func stopServer(server *http.Server, ctx context.Context) error {
	// Заставляем завершиться сервер и ждём его завершения
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Log.Errorf("Failed to shutdown server: %v", err)
	}
	logger.Log.Info("Server stop")

	return err
}
