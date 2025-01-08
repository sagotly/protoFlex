package system_scripts

import (
	"fmt"
)

// SetupNamespace sets up a network namespace and associated configurationsfunc isSubnetUnused(subnet string) (bool, error) {
func SetupNamespace(interfaceName string) error {
	// 0. Detect external interface
	extInterface, err := detectExtInterface()
	if err != nil {
		return fmt.Errorf("failed to detect external interface: %w", err)
	}

	// 1. Find a unique subnet
	subnet, err := findUnusedSubnet()
	if err != nil {
		return fmt.Errorf("failed to find unused subnet: %w", err)
	}
	fmt.Printf("Using subnet: %s\n", subnet)

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
	veth0 := interfaceName + "0"
	veth1 := interfaceName + "1"
	checkVethCmd := fmt.Sprintf("ip link show | grep -w %s", veth0)
	if _, err := runCommand(checkVethCmd); err == nil {
		return fmt.Errorf("veth pair %s and %s already exists, skipping creation.\n", veth0, veth1)
	}
	// Create veth pair
	createVethCmd := fmt.Sprintf("ip link add %s type veth peer name %s", veth0, veth1)
	if _, err := runCommand(createVethCmd); err != nil {
		return fmt.Errorf("failed to create veth pair: %w", err)
	}
	// Assign veth1 to the namespace
	assignVethCmd := fmt.Sprintf("ip link set %s netns %s", veth1, namespace)
	if _, err := runCommand(assignVethCmd); err != nil {
		return fmt.Errorf("failed to assign veth to namespace: %w", err)
	}

	// 4. Configure IP addresses
	hostIP := subnet[:len(subnet)-3] + "1/24"      // E.g., 192.168.x.1/24
	namespaceIP := subnet[:len(subnet)-3] + "2/24" // E.g., 192.168.x.2/24

	configureHostVethCmd := fmt.Sprintf("ip addr add %s dev %s && ip link set %s up", hostIP, veth0, veth0)
	if _, err := runCommand(configureHostVethCmd); err != nil {
		return fmt.Errorf("failed to configure host veth: %w", err)
	}

	configureNamespaceVethCmd := fmt.Sprintf(
		"ip netns exec %s bash -c \"ip addr add %s/24 dev %s && ip link set %s up && ip link set lo up\"",
		namespace, namespaceIP, veth1, veth1,
	)
	if _, err := runCommand(configureNamespaceVethCmd); err != nil {
		return fmt.Errorf("failed to configure namespace veth: %w", err)
	}

	// 5. Set up default route in namespace
	setupRouteCmd := fmt.Sprintf("ip netns exec %s ip route add default via %s dev %s", namespace, hostIP, veth1)
	if _, err := runCommand(setupRouteCmd); err != nil {
		return fmt.Errorf("failed to set default route: %w", err)
	}

	// 6. Configure DNS
	setupDNSCmd := fmt.Sprintf("mkdir -p /etc/netns/%s && echo 'nameserver 1.1.1.1' > /etc/netns/%s/resolv.conf", namespace, namespace)
	if _, err := runCommand(setupDNSCmd); err != nil {
		return fmt.Errorf("failed to configure DNS: %w", err)
	}

	// 7. Set up iptables rules
	setupIptablesCmd := fmt.Sprintf(
		"iptables -A FORWARD -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT &&  iptables -A FORWARD -i %s -j ACCEPT &&  iptables -t nat -A POSTROUTING -s %s -o %s -j MASQUERADE",
		veth0, subnet, extInterface,
	)
	if _, err := runCommand(setupIptablesCmd); err != nil {
		return fmt.Errorf("failed to set up iptables: %w", err)
	}

	// 8. Bring up WireGuard in namespace
	wgCmd := fmt.Sprintf("ip netns exec %s wg-quick up wg0", namespace)
	if _, err := runCommand(wgCmd); err != nil {
		return fmt.Errorf("failed to bring up WireGuard: %w", err)
	}

	fmt.Printf("Namespace %s with veth pair %s, %s, and WireGuard successfully configured using subnet %s.\n", namespace, veth0, veth1, subnet)
	return nil
}
