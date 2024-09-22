package main

import (
	"context"
	"errors"
	"go.uber.org/zap"
	config "gofemart/internal/configuration"
	database "gofemart/internal/databse"
	"gofemart/internal/logger"
	"gofemart/internal/ordercheck"
	"gofemart/internal/router"
	"gofemart/internal/server"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	_, err = initLogger(cnf.LogLevel)
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
	if err = runApplication(cnf); err != nil {
		logger.Log.Error(err)
	}

	logger.Log.Info("End program")
}

// runApplication производим старт приложения
func runApplication(cnf *config.CliConfig) error {
	ctx, cancel := context.WithCancel(context.Background()) // Контекст для правильной остановки синхронизации
	defer func() {
		logger.Log.Info("Cancel context")
		cancel()
	}()

	pool, err := database.NewDB(cnf.DatabaseDSN)
	// Инициализируем базу данных
	if err != nil {
		return err
	}
	// Вызываем функцию закрытия базы данных
	defer pool.Close()
	// Производим миграции базы
	if err = pool.Migrate(); err != nil {
		return err
	}

	ordercheck.CheckPool = ordercheck.NewPool(ordercheck.PoolConfig{
		CTX:             ctx,
		QueueSize:       cnf.QueueSize,
		WorkerCount:     cnf.WorkerCount,
		Pause:           cnf.AccrualSenderPause,
		AccrualURL:      cnf.AccrualSystemAddress,
		DBExecutor:      pool.DBx,
		DBCheckDuration: cnf.DBCheckDuration,
	})
	defer ordercheck.CheckPool.Close()

	wg := new(errgroup.Group)
	serv := server.NewServer(ctx, router.NewRouter(pool, cnf), cnf.Address)
	// Запускаем сервер
	wg.Go(func() error {
		sErr := serv.S.ListenAndServe()
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
	serv.Close()

	// Ожидаем завершения всех горутин перед завершением программы
	if err = wg.Wait(); err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("End Server")
	return nil
}

// initLogger инициализируем логер
func initLogger(logLevel string) (*zap.SugaredLogger, error) {
	lgr, err := logger.New(logLevel)
	if err != nil {
		return nil, err
	}
	logger.Log = lgr

	return lgr, nil
}
