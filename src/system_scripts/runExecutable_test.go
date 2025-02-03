package system_scripts_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/sagotly/protoFlex.git/src/system_scripts"
	"github.com/stretchr/testify/suite"
)

type RunExecutableTestSuite struct {
	suite.Suite
	tempDir   string
	namespace string
	hostIP    string
}

func (suite *RunExecutableTestSuite) SetupSuite() {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "test_executables")
	suite.Require().NoError(err)
	suite.tempDir = tempDir

	// Simulate a namespace (replace with actual setup in real environment)
	suite.namespace = "wg0_namespace"

	// Get the host IP
	hostIPOutput := "2001:fb1:47:b1a6:a286:e769:6cc:d20"
	suite.Require().NoError(err)
	suite.hostIP = strings.TrimSpace(string(hostIPOutput))
}

func (suite *RunExecutableTestSuite) TearDownSuite() {
	// Cleanup temporary files and namespace
	suite.NoError(os.RemoveAll(suite.tempDir))
}

func (suite *RunExecutableTestSuite) createTestFile(name string, hostIP string, executable bool) string {
	content := `#!/bin/bash
actual_ip=$(curl -s https://ifconfig.me)
echo $actual_ip
if [ "` + hostIP + `" == "$actual_ip" ]; then
  echo "FAIL: IP matches host"
  exit 1
else
  echo "PASS: IP is isolated"
fi`
	filePath := suite.tempDir + "/" + name
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.Require().NoError(err)
	if executable {
		suite.NoError(os.Chmod(filePath, 0755))
	}
	return filePath
}

func (suite *RunExecutableTestSuite) TestRunShellScript() {
	path := suite.createTestFile("test.sh", suite.hostIP, true)
	err := system_scripts.RunExecutable(suite.namespace, path, nil)
	suite.NoError(err)
}

func (suite *RunExecutableTestSuite) TestRunBinary() {
	// Path for the binary
	binaryPath := suite.tempDir + "/test_binary"

	// Minimal Go program to compile
	sourceCode := `
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	actualIP, err := http.Get("https://ifconfig.me")
	if err != nil {
		fmt.Println("FAIL: Unable to fetch IP:", err)
		os.Exit(1)
	}

	// Ваше значение IP, с которым нужно сравнить (пример)

	// Проверяем, совпадает ли IP
	if actualIP != nil && actualIP.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(actualIP.Body)
		if err != nil {
			fmt.Println("FAIL: Unable to read response body:", err)
			os.Exit(1)
		}

		actualIPString := string(bodyBytes)
		fmt.Println("Actual IP:", actualIPString)
		if actualIPString == "32" {
			fmt.Println("FAIL: IP matches host")
			os.Exit(1)
		} else {
			fmt.Println("PASS: IP is isolated")
		}
	} else {
		fmt.Println("FAIL: Unable to fetch IP")
		os.Exit(1)
	}
}
`

	// Write the source code to a temporary .go file
	sourcePath := suite.tempDir + "/main.go"
	err := os.WriteFile(sourcePath, []byte(sourceCode), 0644)
	suite.Require().NoError(err)

	// Compile the Go program
	cmd := exec.Command("go", "build", "-o", binaryPath, sourcePath)
	err = cmd.Run()
	suite.Require().NoError(err)

	// Run the binary file
	err = system_scripts.RunExecutable(suite.namespace, binaryPath, nil)
	suite.NoError(err)
}

func (suite *RunExecutableTestSuite) TestRunDockerCompose() {
	content := `version: '3.9'
services:
  test:
    image: alpine/curl
    command: sh -c "apk add --no-cache curl && curl -s https://ifconfig.me && sleep 5"

`

	path := suite.tempDir + "/docker-compose.yml"
	err := os.WriteFile(path, []byte(content), 0644)
	suite.Require().NoError(err)
	err = system_scripts.RunExecutable(suite.namespace, path, nil)
	suite.NoError(err)
}

func (suite *RunExecutableTestSuite) TestRunDockerfile() {
	content := `FROM alpine
CMD sh -c 'actual_ip=$(curl -s https://ifconfig.me) && if [ \"` + suite.hostIP + `\" == \"$actual_ip\" ]; then echo \"FAIL: IP matches host\"; exit 1; else echo \"PASS: IP is isolated\"; fi'`
	path := suite.tempDir + "/Dockerfile"
	err := os.WriteFile(path, []byte(content), 0644)
	suite.Require().NoError(err)
	err = system_scripts.RunExecutable(suite.namespace, path, nil)
	suite.NoError(err)
}

func (suite *RunExecutableTestSuite) TestUnsupportedFileType() {
	path := suite.tempDir + "/unsupported.txt"
	err := os.WriteFile(path, []byte("Just a text file"), 0644)
	suite.Require().NoError(err)
	err = system_scripts.RunExecutable(suite.namespace, path, nil)
	suite.Error(err)
	suite.Contains(err.Error(), "unsupported file type")
}

func TestRunExecutableTestSuite(t *testing.T) {
	suite.Run(t, new(RunExecutableTestSuite))
}
