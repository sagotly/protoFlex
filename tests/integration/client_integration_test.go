package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// TokenResponse описывает структуру ответа от сервера при генерации токена.
type TokenResponse struct {
	Token string `json:"token"`
}

// IntegrationTestSuite описывает интеграционный набор тестов для генерации и валидации токена.
type IntegrationTestSuite struct {
	suite.Suite

	// Статистика для запроса /generate
	generateDuration    time.Duration
	generateStatus      int
	generateResponse    string
	generateRespBytes   int
	generateContentType string
	generateTimestamp   time.Time

	// Статистика для запроса /validate
	validateDuration    time.Duration
	validateStatus      int
	validateResponse    string
	validateRespBytes   int
	validateContentType string
	validateTimestamp   time.Time
}

// TestTokenGenerationAndValidation выполняет запрос к /generate и затем к /validate?token=...
func (s *IntegrationTestSuite) TestTokenGenerationAndValidation() {
	// Отправляем запрос на генерацию токена.
	generateURL := "http://localhost:8088/generate"
	s.generateTimestamp = time.Now()
	startGen := time.Now()
	respGen, err := http.Get(generateURL)
	s.Require().NoError(err, "Ошибка запроса к /generate")
	s.generateDuration = time.Since(startGen)
	s.generateStatus = respGen.StatusCode

	bodyGen, err := ioutil.ReadAll(respGen.Body)
	respGen.Body.Close()
	s.Require().NoError(err, "Ошибка чтения ответа /generate")
	s.generateResponse = string(bodyGen)
	s.generateRespBytes = len(bodyGen)
	s.generateContentType = respGen.Header.Get("Content-Type")

	// Парсим JSON-ответ для извлечения токена.
	var tokenResp TokenResponse
	err = json.Unmarshal(bodyGen, &tokenResp)
	s.Require().NoError(err, "Ошибка парсинга JSON ответа /generate")
	s.Require().NotEmpty(tokenResp.Token, "Токен не получен от /generate")

	// Формируем корректный URL для валидации с чистым значением токена.
	validateURL := fmt.Sprintf("http://localhost:8088/validate?token=%s", tokenResp.Token)
	s.validateTimestamp = time.Now()
	startVal := time.Now()
	respVal, err := http.Get(validateURL)
	s.Require().NoError(err, "Ошибка запроса к /validate")
	s.validateDuration = time.Since(startVal)
	s.validateStatus = respVal.StatusCode

	bodyVal, err := ioutil.ReadAll(respVal.Body)
	respVal.Body.Close()
	s.Require().NoError(err, "Ошибка чтения ответа /validate")
	s.validateResponse = string(bodyVal)
	s.validateRespBytes = len(bodyVal)
	s.validateContentType = respVal.Header.Get("Content-Type")

	// Записываем статистику в текстовый файл.
	file, err := os.Create("integration_stats.txt")
	s.Require().NoError(err, "Не удалось создать файл статистики")
	defer file.Close()

	// Формируем отчёт по статистике.
	fmt.Fprintln(file, "=== Интеграционный тест: Генерация и валидация токена ===")
	fmt.Fprintln(file, "\n[1] Запрос на генерацию токена:")
	fmt.Fprintf(file, "URL: %s\n", generateURL)
	fmt.Fprintf(file, "Время запроса: %d ms\n", s.generateDuration.Milliseconds())
	fmt.Fprintf(file, "Статус код: %d\n", s.generateStatus)
	fmt.Fprintf(file, "Ответ сервера: %s\n", s.generateResponse)
	fmt.Fprintf(file, "Размер ответа: %d байт\n", s.generateRespBytes)
	fmt.Fprintf(file, "Content-Type: %s\n", s.generateContentType)
	fmt.Fprintf(file, "Время начала запроса: %s\n", s.generateTimestamp.Format(time.RFC3339))

	fmt.Fprintln(file, "\n[2] Запрос на валидацию токена:")
	fmt.Fprintf(file, "URL: %s\n", validateURL)
	fmt.Fprintf(file, "Время запроса: %d ms\n", s.validateDuration.Milliseconds())
	fmt.Fprintf(file, "Статус код: %d\n", s.validateStatus)
	fmt.Fprintf(file, "Ответ сервера: %s\n", s.validateResponse)
	fmt.Fprintf(file, "Размер ответа: %d байт\n", s.validateRespBytes)
	fmt.Fprintf(file, "Content-Type: %s\n", s.validateContentType)
	fmt.Fprintf(file, "Время начала запроса: %s\n", s.validateTimestamp.Format(time.RFC3339))

	totalDuration := s.generateDuration + s.validateDuration
	fmt.Fprintln(file, "\n[Общая статистика]")
	fmt.Fprintf(file, "Общее время выполнения обоих запросов: %d ms\n", totalDuration.Milliseconds())
	fmt.Fprintln(file, "\nДополнительная информация:")
	fmt.Fprintln(file, "1. Убедитесь, что сервер запущен на localhost:8088 и корректно обрабатывает запросы.")
	fmt.Fprintln(file, "2. Если валидация токена завершилась неуспешно, проверьте логи сервера или спецификацию API.")
	fmt.Fprintln(file, "3. Обратите внимание на заголовки Content-Type – они должны соответствовать ожидаемому формату (например, application/json).")
	fmt.Fprintln(file, "4. Размер ответа в байтах может помочь отследить, не передаются ли лишние данные.")
	fmt.Fprintln(file, "5. Время начала запроса позволяет соотнести тест с серверными логами и выявить возможные задержки в сети.")
}

// Точка входа для запуска набора интеграционных тестов.
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
