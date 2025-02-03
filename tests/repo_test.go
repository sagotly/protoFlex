package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	enteties "github.com/sagotly/protoFlex.git/src/entities"
	"github.com/sagotly/protoFlex.git/src/repo"
	"github.com/sagotly/protoFlex.git/src/utils"
	"github.com/stretchr/testify/suite"
)

type RepoTestSuite struct {
	suite.Suite
	Db         *sql.DB
	ServerRepo *repo.ServerRepo
	TunnelRepo *repo.TunnelRepo
}

func (suite *RepoTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", "../../test.db")
	if err != nil {
		suite.T().Fatal(err)
	}
	utils.SetupDatabase(db)

	suite.Db = db
	suite.ServerRepo = repo.NewServerRepo(db)
	suite.TunnelRepo = repo.NewTunnelRepo(db)
}

func (suite *RepoTestSuite) TearDownTest() {
	suite.Require().NoError(suite.Db.Close())
	suite.Require().NoError(os.Remove("../../test.db"))
}

func (suite *RepoTestSuite) TestServerAndTunnelRepoIntegration() {
	// Insert a test server
	server := enteties.Server{
		Ip:         "127.0.0.1",
		Name:       "Test Server",
		TunnelList: "[]",
	}
	id, err := suite.ServerRepo.CreateServer(server)
	suite.Require().NoError(err)

	// Fetch server to verify insertion
	fetchedServer, err := suite.ServerRepo.GetServerById(id)
	suite.Require().NoError(err)
	suite.Equal(server.Name, fetchedServer.Name)

	// Fetch all servers
	allServers, err := suite.ServerRepo.GetAllServers()
	suite.Require().NoError(err)
	suite.NotEmpty(allServers)

	// Insert a tunnel linked to the server
	tunnel := enteties.Tunnel{
		ServerId:             fetchedServer.Id,
		InterfaceName:        "wg0",
		ConnectedConnections: `["3434", "1234"]`,
	}
	err = suite.TunnelRepo.CreateTunnel(tunnel)
	suite.Require().NoError(err)

	// Add a connection to the tunnel
	err = suite.TunnelRepo.AddConnectionToTunnel("wg0", "12345")
	suite.Require().NoError(err)

	// Fetch tunnel to verify insertion
	fetchedTunnel, err := suite.TunnelRepo.GetTunnelById(1)
	suite.Require().NoError(err)
	suite.Equal(tunnel.InterfaceName, fetchedTunnel.InterfaceName)

	suite.Contains(fetchedTunnel.ConnectedConnections, "12345")

	// Fetch all tunnels
	allTunnels, err := suite.TunnelRepo.GetAllTunnels()
	suite.Require().NoError(err)
	suite.NotEmpty(allTunnels)

}

func (s *RepoTestSuite) TestGetTunnelByInterfaceName() {
	// Insert a test server
	server := enteties.Server{
		Ip:         "9",
		Name:       "Test Server",
		TunnelList: "[]",
	}
	_, err := s.ServerRepo.CreateServer(server)
	s.Require().NoError(err)
	err = s.TunnelRepo.CreateTunnel(enteties.Tunnel{
		ServerId:             9,
		InterfaceName:        "wggg",
		ConnectedConnections: `[]`,
	})
	s.Require().NoError(err)
	_, exists, err := s.TunnelRepo.GetTunnelByInterfaceName("wggg")

	if err != nil {
		log.Println("Error creating server:", err)
	} else if !exists {
		log.Println("Error creating server: tunnel already exists")
	}
	log.Println(exists)

}
func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}
