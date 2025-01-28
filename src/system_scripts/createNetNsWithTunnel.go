package system_scripts

import (
	"fmt"
	"strings"
)

// SetupNamespace sets up a network namespace and associated configurationsfunc isSubnetUnused(subnet string) (bool, error) {
func SetupNamespace(interfaceName string) error {
	// 0. Detect external interface
	extInterface, err := detectExtInterface()
	if err != nil {
		return fmt.Errorf("failed to detect external interface: %w", err)
	}

	// 1. Find a unique subnet, например "192.168.100.0/24"
	subnetStr, err := findUnusedSubnet()
	if err != nil {
		return fmt.Errorf("failed to find unused subnet: %w", err)
	}
	fmt.Printf("Using subnet: %s\n", subnetStr)

	// Parse the subnet
	hostIPWithMask, nsIPWithMask, err := splitCIDRandMakeHostNS(subnetStr)
	if err != nil {
		return fmt.Errorf("failed to split CIDR and make host namespace: %w", err)
	}
	// 2. Check if namespace exists
	namespace := interfaceName + "_namespace"
	checkNamespaceCmd := fmt.Sprintf("ip netns list | grep -w %s", namespace)
	if _, err := runCommand(checkNamespaceCmd); err == nil {
		fmt.Printf("Namespace %s already exists, skipping creation.\n", namespace)
	} else {
		// Create the namespace
		createNamespaceCmd := fmt.Sprintf("ip netns add %s", namespace)
		if _, err := runCommand(createNamespaceCmd); err != nil {
			return fmt.Errorf("failed to create namespace: %w", err)
		}
	}

	// 3. Check if veth pair exists
	veth0 := "veth_" + interfaceName + "_0"
	veth1 := "veth_" + interfaceName + "_1"
	checkVethCmd := fmt.Sprintf("ip link show | grep -w %s", veth0)
	if _, err := runCommand(checkVethCmd); err == nil {
		fmt.Printf("veth pair %s and %s already exists, skipping creation.\n", veth0, veth1)
	} else {
		// Create veth pair
		createVethCmd := fmt.Sprintf("ip link add %s type veth peer name %s", veth0, veth1)
		if _, err := runCommand(createVethCmd); err != nil {
			return fmt.Errorf("failed to create veth pair: %w", err)
		}
	}
	// Assign veth1 to the namespace
	assignVethCmd := fmt.Sprintf("ip link set %s netns %s", veth1, namespace)
	if _, err := runCommand(assignVethCmd); err != nil {
		return fmt.Errorf("failed to assign veth to namespace: %w", err)
	}

	// 4. Configure IP addresses (host side)
	configureHostVethCmd := fmt.Sprintf("sudo ip addr add %s dev %s && sudo ip link set %s up", hostIPWithMask, veth0, veth0)
	if _, err := runCommand(configureHostVethCmd); err != nil {
		return fmt.Errorf("failed to configure host veth: %w", err)
	}

	// Configure IP addresses (namespace side)
	configureNamespaceVethCmd := fmt.Sprintf(
		"sudo ip netns exec %s bash -c \"ip addr add %s dev %s && sudo ip link set %s up && sudo ip link set lo up\"",
		namespace, nsIPWithMask, veth1, veth1,
	)
	if _, err := runCommand(configureNamespaceVethCmd); err != nil {
		return fmt.Errorf("failed to configure namespace veth: %w", err)
	}

	// 5. Set up default route in namespace (via hostIPRaw)
	setupRouteCmd := fmt.Sprintf("ip netns exec %s ip route add default via %s dev %s", namespace, strings.Split(hostIPWithMask, "/")[0], veth1)
	if _, err := runCommand(setupRouteCmd); err != nil {
		return fmt.Errorf("failed to set default route: %w", err)
	}

	// 6. Configure DNS
	setupDNSCmd := fmt.Sprintf("echo 'nameserver 1.1.1.1' | sudo tee /etc/netns/%s/resolv.conf", namespace)
	if _, err := runCommand(setupDNSCmd); err != nil {
		return fmt.Errorf("failed to configure DNS: %w", err)
	}

	// (Не забывайте про sysctl -w net.ipv4.ip_forward=1, чтобы пакеты могли маршрутизироваться)
	if _, err := runCommand("sudo sysctl -w net.ipv4.ip_forward=1"); err != nil {
		fmt.Printf("Warning: Failed to enable ip_forward: %v\n", err)
	}

	// 7. Set up iptables rules
	// Маска будет та же, ipNet.String() даст, например, "192.168.100.0/24"
	setupIptablesCmd := fmt.Sprintf(
		"sudo iptables -A FORWARD -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT && "+
			"sudo iptables -A FORWARD -i %s -j ACCEPT && "+
			"sudo iptables -t nat -A POSTROUTING -s %s -o %s -j MASQUERADE",
		veth0, subnetStr, extInterface,
	)
	if _, err := runCommand(setupIptablesCmd); err != nil {
		return fmt.Errorf("failed to set up iptables: %w", err)
	}

	// 8. Удаляем старый WireGuard-интерфейс (на случай, если он уже существовал)
	delWGCmd := fmt.Sprintf("sudo ip netns exec %s ip link del %s 2>/dev/null", namespace, interfaceName)
	if _, err := runCommand(delWGCmd); err != nil {
		fmt.Printf("Note: Attempt to delete existing %s returned: %v\n", interfaceName, err)
	}

	// 9. Поднимаем WireGuard-интерфейс с именем interfaceName
	wgCmd := fmt.Sprintf("sudo ip netns exec %s wg-quick up %s", namespace, interfaceName)
	if _, err := runCommand(wgCmd); err != nil {
		return fmt.Errorf("failed to bring up WireGuard: %w", err)
	}

	fmt.Printf("Namespace %s with veth pair %s/%s and WireGuard (%s) successfully configured using subnet %s.\n",
		namespace, veth0, veth1, interfaceName, subnetStr)

	return nil
}

func splitCIDRandMakeHostNS(cidr string) (string, string, error) {
	// Разделяем строку на "192.168.100.0" и "24"
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("неправильный формат CIDR: %s", cidr)
	}
	ipPart := parts[0]   // "192.168.100.0"
	maskPart := parts[1] // "24"

	// Разделяем IP на октеты ["192","168","100","0"]
	octets := strings.Split(ipPart, ".")
	if len(octets) != 4 {
		return "", "", fmt.Errorf("не IPv4-адрес: %s", ipPart)
	}

	// Для Host IP: делаем копию и ставим последний октет = "1"
	hostOctets := make([]string, 4)
	copy(hostOctets, octets)
	hostOctets[3] = "1"
	hostIP := strings.Join(hostOctets, ".") + "/" + maskPart

	// Для Namespace IP: делаем копию и ставим последний октет = "2"
	nsOctets := make([]string, 4)
	copy(nsOctets, octets)
	nsOctets[3] = "2"
	nsIP := strings.Join(nsOctets, ".") + "/" + maskPart

	return hostIP, nsIP, nil
}
