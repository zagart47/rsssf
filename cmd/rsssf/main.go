package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"rsssf/config"
	"rsssf/router"
	"rsssf/service"
	"rsssf/storage"
	"rsssf/storage/postgres"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dsn := config.Configs.Postgres.DSN
	// создаем подключение к БД
	db := postgres.New(dsn)
	// запускаем миграцию
	m, err := migrate.New(
		"file://././storage/postgres/migrations",
		dsn,
	)
	if err != nil {
		log.Println("migrations error:", err.Error(), "!")
	}
	if err = m.Up(); err != nil {
		log.Println("migrations error:", err.Error(), "!")
	}
	// создаем сторедж
	storages := storage.NewStorages(db)
	// создаем сервисы
	services := service.NewServices(storages)
	// создаем роутер
	r := router.NewRouter(services)
	// инициализируем хендлеры
	r.InitHandlers()
	// запускам обновление новостей
	go func() {
		r.UpdateNews()
	}()
	// запускаем сервер
	go func() {
		err := http.ListenAndServe("localhost:8080", &r)
		if err != nil {
			return
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit // Блокируем выполнение до получения сигнала
	log.Println("Shutting down server...")

}
