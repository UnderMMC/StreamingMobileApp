package app

import (
	"StreamingMobileApp/internal/domain/repository"
	"StreamingMobileApp/internal/domain/service"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Service interface {
	GettingStreamFrames(frame []byte) error
}

type App struct {
	serv Service
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем подключения с любого источника
	},
}

func (a *App) GetStream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading websocket", err)
	}

	// Бесконечный цикл для получения кадров видео
	for {
		// Чтение данных (кадра видео)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message", err)
			break
		}

		err = a.serv.GettingStreamFrames(msg)
		if err != nil {
			log.Println("Error saving frame", err)
		}
	}
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	// Инициализация зависимостей слоев
	repo := repository.NewRepository()
	serv := service.NewService(repo)
	a.serv = serv

	http.HandleFunc("/stream", a.GetStream)

	// Запуск сервера на порту 8080
	log.Println("Starting server on port: 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
