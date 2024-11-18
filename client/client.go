package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"
	"image/jpeg"
	"log"
	"time"
)

func main() {
	// Открываем камеру
	videoCapture, err := gocv.OpenVideoCapture(0) // 0 - это первое доступное устройство
	if err != nil {
		log.Fatalf("Error opening video capture device: %v", err)
	}
	defer videoCapture.Close()

	// Создаем окно для отображения видео
	window := gocv.NewWindow("Video Stream")
	defer window.Close()

	// Создаем вебсокет клиент
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/stream", nil)
	if err != nil {
		log.Fatalf("Error connecting to WebSocket server: %v", err)
	}
	defer ws.Close()

	// Буфер для хранения изображения
	img := gocv.NewMat()
	defer img.Close()

	for {
		// Чтение кадра с камеры
		if ok := videoCapture.Read(&img); !ok {
			log.Println("Error reading frame from camera")
			break
		}

		// Отображаем кадр на экране (опционально)
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}

		// Преобразуем кадр в формат JPEG
		// img.ToImage() возвращает два значения: image.Image и ошибку
		imageData, err := img.ToImage()
		if err != nil {
			log.Printf("Error converting to image: %v", err)
			continue
		}

		// Создаем буфер для кодирования в JPEG
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, imageData, nil)
		if err != nil {
			log.Printf("Error encoding frame to JPEG: %v", err)
			continue
		}

		// Отправляем кадр в вебсокет
		err = ws.WriteMessage(websocket.BinaryMessage, buf.Bytes())
		if err != nil {
			log.Printf("Error sending frame over WebSocket: %v", err)
			continue
		}

		// Задержка для контроля частоты отправки кадров
		time.Sleep(100 * time.Millisecond) // 10 кадров в секунду (или можно регулировать)
	}

	fmt.Println("Closing application")
}
