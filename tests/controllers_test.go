package tests

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sagotly/protoFlex.git/src/controllers"
	"github.com/sagotly/protoFlex.git/src/repo"
	"github.com/sagotly/protoFlex.git/src/utils"
)

type ServerViewControllerTestSuite struct {
	suite.Suite

	startTime time.Time

	Db                   *sql.DB
	ServerRepo           *repo.ServerRepo
	TunnelRepo           *repo.TunnelRepo
	AddedExecutablesRepo *repo.AddedExecutablesRepo

	ServerViewController *controllers.ServerViewController
	AddedExecCont        *controllers.AddedExecutablesController
}

func (suite *ServerViewControllerTestSuite) SetupSuite() {
	suite.startTime = time.Now()
}

func (suite *ServerViewControllerTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", "../../test.db")
	fmt.Println("Running as user:", os.Getenv("USER"))
	if err != nil {
		suite.T().Fatal(err)
	}
	utils.SetupDatabase(db)

	suite.Db = db
	suite.ServerRepo = repo.NewServerRepo(db)
	suite.TunnelRepo = repo.NewTunnelRepo(db)
	suite.AddedExecutablesRepo = repo.NewAddedExecutablesRepo(db)

	suite.ServerViewController = controllers.NewServerViewController(suite.TunnelRepo, suite.ServerRepo)
	suite.AddedExecCont = controllers.NewAddedExcecutablesController(suite.TunnelRepo, suite.ServerRepo, suite.AddedExecutablesRepo)
}

func (suite *ServerViewControllerTestSuite) TearDownTest() {
	suite.Require().NoError(suite.Db.Close())
	suite.Require().NoError(os.Remove("../../test.db"))
}

func (suite *ServerViewControllerTestSuite) TearDownSuite() {
	// Вычисляем время выполнения тестов
	duration := time.Since(suite.startTime)

	// Читаем статистику памяти
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Дополнительные метрики:
	// 1. Общее количество выделенной памяти (TotalAlloc)
	// 2. Общая системная память, выделенная для Go (Sys)
	// 3. Количество вызовов malloc (Mallocs)
	// 4. Количество вызовов free (Frees)
	// 5. Количество сборок мусора (NumGC)
	// 6. Число запущенных горутин
	numGoroutine := runtime.NumGoroutine()

	// Создаем файл для записи статистики, например "test_stats.txt"
	file, err := os.Create("test_stats.txt")
	if err != nil {
		suite.T().Fatalf("Не удалось создать файл статистики: %v", err)
	}
	defer file.Close()

	// Формируем содержимое файла с 8 колонками:
	// 1. Время выполнения тестов
	// 2. Использовано памяти (Alloc)
	// 3. Общее количество выделенной памяти (TotalAlloc)
	// 4. Системная память (Sys)
	// 5. Количество вызовов malloc (Mallocs)
	// 6. Количество вызовов free (Frees)
	// 7. Число сборок мусора (NumGC)
	// 8. Число запущенных горутин
	content := fmt.Sprintf(
		"Время выполнения тестов: %s\n"+
			"Использовано памяти (Alloc): %d байт\n"+
			"Общее количество выделенной памяти (TotalAlloc): %d байт\n"+
			"Системная память (Sys): %d байт\n"+
			"Количество вызовов malloc (Mallocs): %d\n"+
			"Количество вызовов free (Frees): %d\n"+
			"Число сборок мусора (NumGC): %d\n"+
			"Число запущенных горутин: %d\n",
		duration,
		m.Alloc,
		m.TotalAlloc,
		m.Sys,
		m.Mallocs,
		m.Frees,
		m.NumGC,
		numGoroutine,
	)

	_, err = file.WriteString(content)
	if err != nil {
		suite.T().Fatalf("Не удалось записать статистику в файл: %v", err)
	}
}

func (suite *ServerViewControllerTestSuite) TestCreateNewServerBtn() {
	// Define test data
	serverName := "Test Server"
	serverIp := "192.168.1.1"
	interfaceName := "wg0"

	// Call the function being tested
	err := suite.ServerViewController.CreateNewServerBtn(serverName, serverIp, interfaceName)
	suite.Require().NoError(err)

	// Verify that the server was created
	servers, err := suite.ServerRepo.GetAllServers()
	suite.Require().NoError(err)
	suite.Require().Len(servers, 1)
	suite.Equal(serverName, servers[0].Name)
	suite.Equal(serverIp, servers[0].Ip)

	// // Verify that the tunnel was created
	// tunnels, err := suite.TunnelRepo.GetAllTunnels()
	// suite.Require().NoError(err)
	// suite.Require().Len(tunnels, 1)
	// suite.Equal(interfaceName, tunnels[0].InterfaceName)
	// Ensure the tunnel is linked to the correct server
	// suite.Equal(servers[0].Id, tunnels[0].ServerId)
}

func (suite *ServerViewControllerTestSuite) TestAddExecutableBtn() {
	// Define test data for each server
	serverName := "Test Server"
	serverIp := "192.168.1.1"
	interfaceName := "wg0"

	// Create a new server
	err := suite.ServerViewController.CreateNewServerBtn(serverName, serverIp, interfaceName)
	suite.Require().NoError(err)

	// Create a temporary file containing a simple shell script
	tmpFile, err := os.CreateTemp("", "hello-world-%d-*.sh")
	suite.Require().NoError(err)

	scriptContent := "#!/bin/bash\necho Hello World\n"
	_, err = tmpFile.Write([]byte(scriptContent))
	suite.Require().NoError(err)

	// Close the file and make it executable
	err = tmpFile.Close()
	suite.Require().NoError(err)
	err = os.Chmod(tmpFile.Name(), 0755)
	suite.Require().NoError(err)

	executablePath := tmpFile.Name()
	arguments := []string{}

	// Add the executable to the server
	err = suite.AddedExecCont.AddExecutableBtn(executablePath, arguments, interfaceName)
	suite.Require().NoError(err)

	err = suite.AddedExecCont.ClickOnExecutableBtn(1, executablePath, "")
	suite.Require().NoError(err)
	cmdStart := time.Now()
	cmd := exec.Command("ip", "link", "delete", "veth_wg_0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	suite.Require().NoError(err)
	suite.startTime = suite.startTime.Add(time.Since(cmdStart))
}

func (suite *ServerViewControllerTestSuite) TestServerBinaryConnection() {
	// ————— START: засечём начало теста
	testStart := time.Now()

	// —— Arrange
	serverName := "Test Server"
	serverIP := "192.168.1.1"
	interfaceName := "wg0"
	executablePath := "./test_materials/protoflex-server-test"

	suite.Require().NoError(
		suite.ServerViewController.CreateNewServerBtn(serverName, serverIP, interfaceName),
		"failed to add server",
	)
	suite.Require().NoError(
		suite.AddedExecCont.AddExecutableBtn(executablePath, []string{}, interfaceName),
		"failed to register executable",
	)

	// —— Act #1: запускаем в namespace
	suite.Require().NoError(
		suite.AddedExecCont.ClickOnExecutableBtn(1, executablePath, ""),
		"ClickOnExecutableBtn should start the server binary inside namespace",
	)

	// ждём мильсекунду и пытаемся подключиться — тут ожидание падения
	time.Sleep(100 * time.Millisecond)
	nsConnErr := error(nil)
	nsConn, err := net.DialTimeout("tcp", "localhost:8080", 100*time.Millisecond)
	if err != nil {
		nsConnErr = err
	} else {
		nsConn.Close()
	}

	// —— Act #2: запускаем на хосте
	cmd := exec.Command(executablePath)
	suite.Require().NoError(cmd.Start(), "failed to start server binary on host")
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	// ждём, пока порт откроется
	var hostRespBody string
	var hostStatusCode int
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://localhost:8080")
		if err == nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			hostRespBody = string(body)
			hostStatusCode = resp.StatusCode
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	suite.Require().NotZero(hostStatusCode, "server did not respond on host within timeout")

	// ————— END: замеряем продолжительность и память
	duration := time.Since(testStart)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// ————— Записываем всё в файл
	f, err := os.Create("test_server_in_ns_reach_results.txt")
	suite.Require().NoError(err, "cannot create results file")
	defer f.Close()

	fmt.Fprintf(f, "=== TestServerBinaryConnection results ===\n\n")
	// namespace
	if nsConnErr != nil {
		fmt.Fprintf(f, "Namespace: connection error (as expected): %v\n", nsConnErr)
	} else {
		fmt.Fprintf(f, "Namespace: unexpectedly connected!\n")
	}
	// host
	fmt.Fprintf(f, "\nHost:\n")
	fmt.Fprintf(f, "  Status Code: %d\n", hostStatusCode)
	fmt.Fprintf(f, "  Body: %q\n", hostRespBody)

	// metrics
	fmt.Fprintf(f, "\nMetrics:\n")
	fmt.Fprintf(f, "  Test Duration: %s\n", duration)
	fmt.Fprintf(f, "  Memory Alloc: %d bytes\n", m.Alloc)
	fmt.Fprintf(f, "  TotalAlloc: %d bytes\n", m.TotalAlloc)
	fmt.Fprintf(f, "  Sys: %d bytes\n", m.Sys)
	fmt.Fprintf(f, "  Mallocs: %d\n", m.Mallocs)
	fmt.Fprintf(f, "  Frees: %d\n", m.Frees)
	fmt.Fprintf(f, "  NumGC: %d\n", m.NumGC)

	// закрываем соединения, если есть
	if nsConn != nil {
		_ = nsConn.Close()
	}
}

func TestServerViewControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerViewControllerTestSuite))
}
