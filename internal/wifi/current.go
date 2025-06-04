package wifi

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func GetCurrentWiFiProfile() (string, error) {
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`\bSSID\s*:\s*(.+)`)
	matches := re.FindAllStringSubmatch(string(output), -1)
	for _, match := range matches {
		if len(match) > 1 {
			return strings.TrimSpace(match[1]), nil
		}
	}
	return "", fmt.Errorf("SSID not found")
}
