package controllers_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sagotly/protoFlex.git/src/controllers"
	enteties "github.com/sagotly/protoFlex.git/src/entities"
	"github.com/sagotly/protoFlex.git/src/repo"
	"github.com/sagotly/protoFlex.git/src/utils"
)

type ServerViewControllerTestSuite struct {
	suite.Suite
	Db                   *sql.DB
	ServerRepo           *repo.ServerRepo
	TunnelRepo           *repo.TunnelRepo
	AddedExecutablesRepo *repo.AddedExecutablesRepo

	ServerViewController *controllers.ServerViewController
	AddedExecCont        *controllers.AddedExecutablesController
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
	// Define test data
	executablePath := "/usr/bin/test-exec"
	arguments := []string{"--arg1", "--arg2"}
	tunnelInterface := "wg0"

	// Ensure the test environment has at least one tunnel to link the executable
	server := enteties.Server{Ip: "127.0.0.1", Name: "Test Server", TunnelList: "[]"}
	id, err := suite.ServerRepo.CreateServer(server)
	suite.Require().NoError(err)

	tunnel := enteties.Tunnel{ServerId: id, InterfaceName: tunnelInterface}
	err = suite.TunnelRepo.CreateTunnel(tunnel)
	suite.Require().NoError(err)

	// Call the function being tested
	err = suite.AddedExecCont.AddExecutableBtn(executablePath, arguments, tunnelInterface)
	suite.Require().NoError(err)

	// Verify that the executable was added
	executables, err := suite.AddedExecutablesRepo.GetAllAddedExecutabless()
	suite.Require().NoError(err)
	suite.Require().Len(executables, 1)

	// Check the details of the added executable
	// suite.Equal(executablePath, executables[0].Path)
	// suite.Equal(tunnelInterface, executables[0].TunnelId)
	// suite.ElementsMatch(arguments, executables[0].Arguments)
}

func TestServerViewControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerViewControllerTestSuite))
}
