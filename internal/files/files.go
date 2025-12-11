// internal/files/files.go
package files

import (
	"BastetTetlegram/internal/models" // Импортируем модели
	"encoding/json"
	"io/ioutil"
	"log" // Стандартная библиотека
	"os"
	"strings"
	"time"
)

// ReadPhrasesFromFile читает фразы из файла
func ReadPhrasesFromFile(filename string) ([]string, error) {
	log.Printf("Попытка чтения файла фраз: %s", filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Ошибка чтения файла фраз: %v", err)
		return nil, err
	}
	text := string(content)
	lines := strings.Split(text, "\n")
	var phrases []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			phrases = append(phrases, trimmedLine)
		}
	}
	log.Printf("Успешно прочитано %d фраз из файла %s", len(phrases), filename)
	return phrases, nil
}

// ReadToastsFromFile читает тосты из файла
func ReadToastsFromFile(filename string) ([]string, error) {
	log.Printf("Попытка чтения файла тостов: %s", filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Ошибка чтения файла тостов: %v", err)
		return nil, err
	}
	text := string(content)
	parts := strings.Split(text, "* * *")
	var toasts []string
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart != "" {
			toasts = append(toasts, trimmedPart)
		}
	}
	log.Printf("Успешно прочитано %d тостов из файла %s", len(toasts), filename)
	return toasts, nil
}

// LoadLastMentionFromFile загружает время из файла
func LoadLastMentionFromFile(filename string) (time.Time, error) {
	log.Printf("Попытка загрузки времени из файла: %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Файл %s не найден, будет создан при следующем обновлении.", filename)
			return time.Time{}, err
		}
		log.Printf("Ошибка открытия файла: %v", err)
		return time.Time{}, err
	}
	defer file.Close()

	var data models.LastMentionData // Используем структуру из models
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		log.Printf("Ошибка декодирования JSON из файла: %v", err)
		return time.Time{}, err
	}

	log.Printf("Время успешно загружено из файла: %v", data.LastMention)
	return data.LastMention, nil
}

// SaveLastMentionToFile сохраняет время в файл
func SaveLastMentionToFile(filename string, lastMention time.Time) error {
	log.Printf("Сохранение времени в файл: %s, время: %v", filename, lastMention)
	data := models.LastMentionData{LastMention: lastMention} // Используем структуру из models

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Ошибка создания файла для сохранения: %v", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		log.Printf("Ошибка кодирования JSON для сохранения: %v", err)
		return err
	}

	log.Printf("Время успешно сохранено в файл.")
	return nil
}
