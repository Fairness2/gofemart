package application

import (
	"context"
	"errors"
	config "gofemart/internal/configuration"
	database "gofemart/internal/databse"
	"gofemart/internal/logger"
	"gofemart/internal/ordercheck"
	"gofemart/internal/router"
	"gofemart/internal/server"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// New производим старт приложения
func New(cnf *config.CliConfig) error {
	ctx, cancel := context.WithCancel(context.Background()) // Контекст для правильной остановки синхронизации
	defer func() {
		logger.Log.Info("Cancel context")
		cancel()
	}()

	pool, err := database.NewDB(cnf.DatabaseDSN, cnf.DBMaxConnections, cnf.DBMaxIdleConnections)
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
