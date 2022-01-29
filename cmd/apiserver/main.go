package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/ziyadovea/user-balance/internal/app/apiserver"
	"github.com/ziyadovea/user-balance/internal/app/apiserver/handler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Создаем экземпляр хендлера
	h := handler.NewHandler(nil)

	// Создаем экземпляр сервера
	server := apiserver.NewServer("8080", h.InitRoutes())

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
