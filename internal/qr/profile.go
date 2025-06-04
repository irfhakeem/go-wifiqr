package qr

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func GetProfileDetail(p string) map[string]string {
	data := make(map[string]string)
	cmd := exec.Command("netsh", "wlan", "show", "profile", p, "key=clear")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting WiFi profile: %v\n", err)
		return data
	}

	outputString := string(output)
	reAuth := regexp.MustCompile(`Authentication\s*:\s*(.+)`)
	rePw := regexp.MustCompile(`Key Content\s*:\s*(.+)`)

	if match := reAuth.FindStringSubmatch(outputString); len(match) > 1 {
		auth := strings.TrimSpace(match[1])
		if strings.Contains(auth, "WPA") {
			auth = "WPA"
		} else if strings.Contains(auth, "WEP") {
			auth = "WEP"
		}
		data["authentication"] = auth
	}

	if match := rePw.FindStringSubmatch(outputString); len(match) > 1 {
		data["password"] = strings.TrimSpace(match[1])
	}

	return data
}

func ParseProfiles(output string) []string {
	var profiles []string
	re := regexp.MustCompile(`All User Profile\s*:\s*(.+)`)
	matches := re.FindAllStringSubmatch(output, -1)
	for _, match := range matches {
		if len(match) > 1 {
			profiles = append(profiles, strings.TrimSpace(match[1]))
		}
	}
	return profiles
}
