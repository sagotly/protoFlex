package ui_test

// import (
// 	"testing"

// 	"database/sql"

// 	"fyne.io/fyne/v2/app"
// 	_ "github.com/mattn/go-sqlite3"
// 	"github.com/sagotly/protoFlex.git/src/controllers"
// 	"github.com/sagotly/protoFlex.git/src/repo"
// 	"github.com/sagotly/protoFlex.git/src/ui"
// 	Ui "github.com/sagotly/protoFlex.git/src/ui"
// 	"github.com/sagotly/protoFlex.git/src/utils"
// 	"github.com/stretchr/testify/suite"
// )

// type UISuite struct {
// 	suite.Suite
// 	Db                   *sql.DB
// 	ServerRepo           *repo.ServerRepo
// 	TunnelRepo           *repo.TunnelRepo
// 	addedExecutablesRepo *repo.AddedExecutablesRepo
// 	UIInstance           *ui.Ui
// }

// func (suite *UISuite) SetupTest() {
// 	var err error
// 	suite.Db, err = sql.Open("sqlite3", "../../test.db")
// 	suite.Require().NoError(err)

// 	suite.Require().NoError(utils.SetupDatabase(suite.Db))

// 	suite.ServerRepo = repo.NewServerRepo(suite.Db)
// 	suite.TunnelRepo = repo.NewTunnelRepo(suite.Db)
// 	suite.addedExecutablesRepo = repo.NewAddedExecutablesRepo(suite.Db)

// 	appInstance := app.New()
// 	window := appInstance.NewWindow("Protoflex Test")

// 	addedExecutablesController := controllers.NewAddedExcecutablesController(suite.TunnelRepo, suite.ServerRepo, suite.addedExecutablesRepo)
// 	serverViewConntroller := controllers.NewServerViewController(suite.TunnelRepo, suite.ServerRepo)
// 	suite.UIInstance = Ui.NewUI(window, serverViewConntroller, addedExecutablesController, suite.ServerRepo, suite.TunnelRepo, suite.addedExecutablesRepo)
// }

// func (suite *UISuite) TearDownTest() {
// 	suite.Require().NoError(suite.Db.Close())
// }

// func (suite *UISuite) TestBuildServerView() {
// 	container, err := suite.UIInstance.BuildUi()
// 	suite.Require().NoError(err)
// 	suite.NotNil(container)
// }

// func (suite *UISuite) TestBuildStatisticsView() {
// 	container, err := suite.UIInstance.BuildUi()
// 	suite.Require().NoError(err)
// 	suite.NotNil(container)
// }

// func (suite *UISuite) TestBuildActiveConnectionsView() {
// 	container, err := suite.UIInstance.BuildUi()
// 	suite.Require().NoError(err)
// 	suite.NotNil(container)
// }

// func TestUISuite(t *testing.T) {
// 	suite.Run(t, new(UISuite))
// }
