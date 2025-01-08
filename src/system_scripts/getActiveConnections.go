package system_scripts

import (
	"fmt"
	"strings"

	enteties "github.com/sagotly/protoFlex.git/src/entities"
)

func GetActiveConnections() (enteties.ActiveConnections, error) {
	// Execute the netstat command
	cmdOutput, err := runCommand("netstat -tunp")
	if err != nil {
		return nil, fmt.Errorf("error running netstat: %w", err)
	}

	lines := strings.Split(cmdOutput, "\n")
	connections := enteties.ActiveConnections{}

	// Skip the header and parse each line
	for _, line := range lines {
		if strings.Contains(line, "tcp") || strings.Contains(line, "udp") {
			fields := strings.Fields(line)
			if len(fields) < 7 {
				continue // Skip invalid or malformed entries
			}

			// Parse the protocol (TCP or UDP)
			protocol := fields[0]

			// Extracting the source and destination addresses
			sourceAddress := fields[3] // Source IP:Port
			destAddress := fields[4]   // Destination IP:Port

			// Splitting IP and port for source and destination
			sourceParts := strings.Split(sourceAddress, ":")
			destParts := strings.Split(destAddress, ":")
			if len(sourceParts) < 2 || len(destParts) < 2 {
				continue // Malformed IP:Port
			}

			sourceIP, sourcePort := sourceParts[0], sourceParts[1]
			destIP, destPort := destParts[0], destParts[1]

			// Parse the PID and process name from the last field
			pidProcess := fields[6]
			pid := strings.Split(pidProcess, "/")[0]

			// Generate 5-tuple string representation
			tuple := fmt.Sprintf("%s:%s@%s:%s %s", sourceIP, sourcePort, destIP, destPort, protocol)

			connections = append(connections, &enteties.ActiveConnection{
				Pid:              pid,
				TypeOfConnection: protocol,
				FiveTuple:        tuple,
			})
		}
	}
	return connections, nil
}
