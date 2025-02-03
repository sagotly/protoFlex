package tests

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sagotly/protoFlex.git/src/client"
	"github.com/stretchr/testify/suite"
)

// ServerClientTestSuite определяет набор тестов для ServerClient,
// а также собирает статистику по запросам.
type ServerClientTestSuite struct {
	suite.Suite
	server *httptest.Server
	client *client.ServerClient
	host   string
	port   string

	// Статистика, собираемая через обработчик тестового сервера:
	generateCount int32
	validateCount int32

	// Дополнительная статистика времени (в миллисекундах) для вызовов:
	totalGenerateDurationMs      int64
	totalValidateDurationMs      int64
	generateErrorCount           int
	validateErrorCount           int
	totalGenerateErrorDurationMs int64
	totalValidateErrorDurationMs int64
}

// SetupTest подготавливает тестовую среду перед каждым тестом.
// Создается тестовый HTTP-сервер, который отвечает на запросы к /generate и /validate,
// обновляя счетчики статистики.
func (s *ServerClientTestSuite) SetupTest() {
	// Обнуляем счетчики статистики
	s.generateCount = 0
	s.validateCount = 0
	s.totalGenerateDurationMs = 0
	s.totalValidateDurationMs = 0
	s.generateErrorCount = 0
	s.validateErrorCount = 0
	s.totalGenerateErrorDurationMs = 0
	s.totalValidateErrorDurationMs = 0

	// Создаем сервер с обработчиком для /generate и /validate
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/generate":
			atomic.AddInt32(&s.generateCount, 1)
			w.WriteHeader(http.StatusOK)
		case "/validate":
			atomic.AddInt32(&s.validateCount, 1)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	// Извлекаем host и port из адреса сервера
	addr := s.server.Listener.Addr().String()
	var err error
	s.host, s.port, err = net.SplitHostPort(addr)
	s.Require().NoError(err, "Не удалось разделить host и port")
	s.client = client.NewServerClient()
}

// TearDownTest освобождает ресурсы после каждого теста.
func (s *ServerClientTestSuite) TearDownTest() {
	s.server.Close()
}

// TearDownSuite вызывается после выполнения всех тестов.
// Здесь генерируется текстовый файл со статистикой запросов в миллисекундах.
func (s *ServerClientTestSuite) TearDownSuite() {
	f, err := os.Create("stats.txt")
	s.Require().NoError(err, "Не удалось создать файл статистики")
	defer f.Close()

	// Статистика для /generate
	fmt.Fprintln(f, "=== Статистика для /generate ===")
	fmt.Fprintf(f, "Количество запросов: %d\n", s.generateCount)
	fmt.Fprintf(f, "Общее время выполнения: %d ms\n", s.totalGenerateDurationMs)
	if s.generateCount > 0 {
		avg := s.totalGenerateDurationMs / int64(s.generateCount)
		fmt.Fprintf(f, "Среднее время запроса: %d ms\n", avg)
	}
	fmt.Fprintf(f, "Количество ошибочных запросов: %d\n", s.generateErrorCount)
	if s.generateErrorCount > 0 {
		avgErr := s.totalGenerateErrorDurationMs / int64(s.generateErrorCount)
		fmt.Fprintf(f, "Среднее время ошибочного запроса: %d ms\n", avgErr)
	}

	// Статистика для /validate
	fmt.Fprintln(f, "\n=== Статистика для /validate ===")
	fmt.Fprintf(f, "Количество запросов: %d\n", s.validateCount)
	fmt.Fprintf(f, "Общее время выполнения: %d ms\n", s.totalValidateDurationMs)
	if s.validateCount > 0 {
		avg := s.totalValidateDurationMs / int64(s.validateCount)
		fmt.Fprintf(f, "Среднее время запроса: %d ms\n", avg)
	}
	fmt.Fprintf(f, "Количество ошибочных запросов: %d\n", s.validateErrorCount)
	if s.validateErrorCount > 0 {
		avgErr := s.totalValidateErrorDurationMs / int64(s.validateErrorCount)
		fmt.Fprintf(f, "Среднее время ошибочного запроса: %d ms\n", avgErr)
	}
}

// callGenerateToken оборачивает вызов GenerateToken с замером времени в миллисекундах.
func (s *ServerClientTestSuite) callGenerateToken(host, port string) (err error, durationMs int64) {
	start := time.Now()
	err = s.client.GenerateToken(host, port)
	duration := time.Since(start)
	durationMs = duration.Milliseconds()
	s.totalGenerateDurationMs += durationMs
	if err != nil {
		s.generateErrorCount++
		s.totalGenerateErrorDurationMs += durationMs
	}
	return err, durationMs
}

// callValidateToken оборачивает вызов ValidateToken с замером времени в миллисекундах.
func (s *ServerClientTestSuite) callValidateToken(host, port, token string) (err error, durationMs int64) {
	start := time.Now()
	err = s.client.ValidateToken(host, port, token)
	duration := time.Since(start)
	durationMs = duration.Milliseconds()
	s.totalValidateDurationMs += durationMs
	if err != nil {
		s.validateErrorCount++
		s.totalValidateErrorDurationMs += durationMs
	}
	return err, durationMs
}

// TestGenerateTokenSuccess проверяет успешный сценарий для GenerateToken.
func (s *ServerClientTestSuite) TestGenerateTokenSuccess() {
	err, durationMs := s.callGenerateToken(s.host, s.port)
	s.Require().NoError(err, "Ожидали успешный ответ при генерации токена")
	s.T().Logf("GenerateToken call duration: %d ms", durationMs)
	s.Equal(int32(1), s.generateCount, "Ожидали 1 запрос к /generate")
}

// TestGenerateTokenFailure проверяет, что функция корректно обрабатывает ошибочный ответ сервера.
func (s *ServerClientTestSuite) TestGenerateTokenFailure() {
	// Создаем временный сервер, который возвращает ошибку для /generate
	failureServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/generate" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer failureServer.Close()

	addr := failureServer.Listener.Addr().String()
	host, port, err := net.SplitHostPort(addr)
	s.Require().NoError(err, "Не удалось разделить host и port для failureServer")

	start := time.Now()
	err = s.client.GenerateToken(host, port)
	durationMs := time.Since(start).Milliseconds()
	s.T().Logf("GenerateToken (failure) call duration: %d ms", durationMs)
	s.Error(err, "Ожидали ошибку, когда сервер возвращает Internal Server Error")
}

// TestValidateTokenSuccess проверяет успешное выполнение ValidateToken.
func (s *ServerClientTestSuite) TestValidateTokenSuccess() {
	err, durationMs := s.callValidateToken(s.host, s.port, "sample-token")
	s.Require().NoError(err, "Ожидали успешную валидацию токена")
	s.T().Logf("ValidateToken call duration: %d ms", durationMs)
	s.Equal(int32(1), s.validateCount, "Ожидали 1 запрос к /validate")
}

// TestValidateTokenFailure проверяет, что ValidateToken корректно обрабатывает ошибочный ответ сервера.
func (s *ServerClientTestSuite) TestValidateTokenFailure() {
	// Создаем временный сервер, который возвращает ошибку для /validate
	failureServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/validate" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer failureServer.Close()

	addr := failureServer.Listener.Addr().String()
	host, port, err := net.SplitHostPort(addr)
	s.Require().NoError(err, "Не удалось разделить host и port для failureServer")

	start := time.Now()
	err = s.client.ValidateToken(host, port, "sample-token")
	durationMs := time.Since(start).Milliseconds()
	s.T().Logf("ValidateToken (failure) call duration: %d ms", durationMs)
	s.Error(err, "Ожидали ошибку при неуспешной валидации токена")
}

// TestRequestStatistics демонстрирует сбор статистики, выполняя несколько вызовов функций GenerateToken и ValidateToken.
func (s *ServerClientTestSuite) TestRequestStatistics() {
	// Выполняем 5 вызовов GenerateToken
	for i := 0; i < 5; i++ {
		err, d := s.callGenerateToken(s.host, s.port)
		s.Require().NoError(err, fmt.Sprintf("GenerateToken (итерация %d) завершился с ошибкой", i))
		s.T().Logf("GenerateToken (итерация %d) duration: %d ms", i, d)
	}

	// Выполняем 3 вызова ValidateToken с разными токенами
	for i := 0; i < 3; i++ {
		token := "token-" + strconv.Itoa(i)
		err, d := s.callValidateToken(s.host, s.port, token)
		s.Require().NoError(err, fmt.Sprintf("ValidateToken (итерация %d) завершился с ошибкой", i))
		s.T().Logf("ValidateToken (итерация %d) duration: %d ms", i, d)
	}

	s.T().Logf("=== Статистика запросов ===")
	s.T().Logf("Количество запросов к /generate: %d", s.generateCount)
	s.T().Logf("Количество запросов к /validate: %d", s.validateCount)
}

// Точка входа для запуска набора тестов.
func TestServerClientTestSuite(t *testing.T) {
	suite.Run(t, new(ServerClientTestSuite))
}
