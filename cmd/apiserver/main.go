package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/ziyadovea/user-balance/config"
	"github.com/ziyadovea/user-balance/internal/app/apiserver"
	"github.com/ziyadovea/user-balance/internal/app/apiserver/handler"
	"github.com/ziyadovea/user-balance/internal/app/repository"
	"github.com/ziyadovea/user-balance/internal/app/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Читаем конфиг
	if err := config.InitConfig(); err != nil {
		log.Fatal().Err(err).Msg("error occurred while reading configuration file")
	}

	// Читаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("error occurred while reading .env file")
	}

	// Создание экземпляра базы данных
	db, err := repository.NewPostgresDB(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		DBName:   viper.GetString("db.db_name"),
		SSLMode:  viper.GetString("db.ssl_mode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error occurred while connecting to the postgres database")
	}

	// Создание репозитория на основе экземпляра БД
	repo := repository.NewRepository(db)

	// Создание сервиса на основе экземпляра репозитория
	services := service.NewService(repo)

	// Создаем экземпляр хендлера
	h := handler.NewHandler(services)

	// Создаем экземпляр сервера
	server := apiserver.NewServer(viper.GetString("port"), h.InitRoutes())

	// Gracefully shutdown
	// Запуск сервера в отдельной горутине
	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("error occurred while running http server")
		}
	}()

	log.Print("Server started")

	quit := make(chan os.Signal, 1)                                    // Создаем буферизированный канал сигналов
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // Указываем, какие UNIX-сигналы будут попадать в наш канал
	<-quit                                                             // Ожидаем сигнала

	// Остановка работы сервера
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("error occurred on server shutting down")
	}

	log.Print("Server stopped")

}
