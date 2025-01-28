package system_scripts

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// runExecutable executes a given file or command based on its type in a specified network namespace.
// Supports scripts, binaries, Docker Compose, and Dockerfile.
func RunExecutable(namespace string, path string, args []string) error {
	switch {
	case isShellScript(path):
		return runShellScript(namespace, path, args)
	case isBinary(path):
		return runBinary(namespace, path, args)
	// case strings.HasSuffix(path, "docker-compose.yml"):
	// 	return runDockerCompose(namespace, path, args)
	// case strings.HasSuffix(path, "Dockerfile"):
	// 	return runDockerfile(namespace, path, args)
	default:
		return errors.New("unsupported file type")
	}
}

// runShellScript executes a shell script with the provided arguments in a namespace.
func runShellScript(namespace string, path string, args []string) error {
	if !hasExecutePermission(path) {
		return errors.New("file does not have execute permission")
	}
	cmd := exec.Command("sudo", append([]string{"ip", "netns", "exec", namespace, "bash", path}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runBinary executes a binary with the provided arguments in a namespace.
func runBinary(namespace string, path string, args []string) error {
	if !hasExecutePermission(path) {
		return errors.New("file does not have execute permission")
	}
	cmd := exec.Command("sudo", append([]string{"ip", "netns", "exec", namespace, path}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runDockerCompose executes Docker Compose with the specified YAML file in a namespace.
// func runDockerCompose(namespace string, path string, args []string) error {
// 	cmd := exec.Command("sudo", append([]string{"ip", "netns", "exec", namespace, "docker-compose", "-f", path, "up"}, args...)...)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	return cmd.Run()
// }

// // runDockerfile builds and runs a container from a Dockerfile in a namespace.
// func runDockerfile(namespace string, path string, args []string) error {
// 	dir := strings.TrimSuffix(path, "/Dockerfile")

// 	// Check if the image already exists
// 	imageName := "temp_image"
// 	checkImageCmd := exec.Command("sudo", "ip", "netns", "exec", namespace, "docker", "images", "-q", imageName)
// 	output, err := checkImageCmd.Output()
// 	if err != nil {
// 		return fmt.Errorf("failed to check existing Docker image: %w", err)
// 	}

// 	if len(strings.TrimSpace(string(output))) == 0 {
// 		// Build the Docker image if it doesn't exist
// 		buildCmd := exec.Command("sudo", append([]string{"ip", "netns", "exec", namespace, "docker", "build", "-t", imageName, dir}, args...)...)
// 		buildCmd.Stdout = os.Stdout
// 		buildCmd.Stderr = os.Stderr
// 		if err := buildCmd.Run(); err != nil {
// 			return fmt.Errorf("failed to build Dockerfile: %w", err)
// 		}
// 	}

// 	// Run the Docker container
// 	runCmd := exec.Command("sudo", append([]string{"ip", "netns", "exec", namespace, "docker", "run", "--rm", imageName}, args...)...)
// 	runCmd.Stdout = os.Stdout
// 	runCmd.Stderr = os.Stderr
// 	return runCmd.Run()
// }

// isBinary checks if a file is an executable binary.
func isBinary(path string) bool {
	if !hasExecutePermission(path) {
		return false
	}
	cmd := exec.Command("file", "--mime-type", "-b", path)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "application/x-executable") ||
		strings.Contains(string(output), "application/x-elf")
}

// isShellScript checks if a file is a shell script.
func isShellScript(path string) bool {
	if !strings.HasSuffix(path, ".sh") {
		return false
	}
	if !hasExecutePermission(path) {
		return false
	}
	return true
}

// hasExecutePermission checks if a file has execute permissions.
func hasExecutePermission(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return info.Mode().Perm()&0111 != 0
}
