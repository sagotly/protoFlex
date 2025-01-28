package system_scripts

import (
	"fmt"
	"os/exec"
	"strings"
)

func isSubnetUnused(subnet string) (bool, error) {
	checkSubnetCmd := fmt.Sprintf("ip route show %s", subnet)
	output, err := runCommand(checkSubnetCmd)
	if err != nil {
		return false, fmt.Errorf("failed to run command: %w", err)
	}
	// Если вывод пустой, значит подсеть свободна
	if strings.TrimSpace(output) == "" {
		return true, nil
	}
	return false, nil
}

func findUnusedSubnet() (string, error) {
	// Test subnet ranges from 192.168.x.0/24 (x from 1 to 254) to find an unused subnet
	for i := 2; i < 255; i++ {
		subnet := fmt.Sprintf("192.168.%d.0/24", i)
		// Check if subnet is unused
		available, err := isSubnetUnused(subnet)
		if err != nil {
			return "", fmt.Errorf("error checking subnet: %w", err)
		}
		if available {
			return subnet, nil
		}
	}
	return "", fmt.Errorf("no unused subnets found")
}

func detectExtInterface() (string, error) {
	// Run a command to list the default route and extract the interface name
	cmd := "ip route | grep default | awk '{print $5}'"
	output, err := runCommand(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to detect external interface: %w", err)
	}

	// Clean and return the result
	return strings.TrimSpace(output), nil
}

func runCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", "sudo "+command)
	fmt.Println("Running command: ", cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %s, %w", string(output), err)
	}
	return string(output), nil
}
