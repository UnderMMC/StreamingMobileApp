package repository

import (
	"fmt"
	"os"
	"time"
)

type Repository struct {
	directory string
}

func NewRepository() *Repository {
	// Указываем директорию для сохранения файлов
	dir := "./frames"
	// Проверка, существует ли директория, если нет — создаем ее
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("Failed to create directory: %v", err))
		}
	}
	return &Repository{directory: dir}

}

func (r *Repository) FrameSaving(frame []byte) error {
	// Генерация имени файла с текущей меткой времени
	filename := fmt.Sprintf("%s/frame_%d.jpg", r.directory, time.Now().UnixNano())

	// Открытие или создание файла
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Сохраняем кадр в файл
	_, err = file.Write(frame)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}
