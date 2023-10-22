package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
	db     *sql.DB
	config Config
}

func New(config Config) *App {
	app := &App{
		rdb:    redis.NewClient(&redis.Options{Addr: config.RedisAddress}),
		config: config,
	}

	//Starts and Configures SQL Database
	db, err := app.StartSqlDatabase()
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	//Configures Tables
	err = createTables(db)
	if err != nil {
		panic(err)
	}

	app.db = db

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	//Checks connections to Redis
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting Server")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}

}
