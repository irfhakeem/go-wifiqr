package wifi

import (
	"fmt"
	"go-wifi-qr/internal/qr"
	"os/exec"
)

func ListAndSelectProfile() {
	cmd := exec.Command("netsh", "wlan", "show", "profiles")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting WiFi profiles: %v\n", err)
		return
	}

	profiles := qr.ParseProfiles(string(output))
	if len(profiles) == 0 {
		fmt.Println("No WiFi profiles found.")
		return
	}

	fmt.Println("Available WiFi Profiles:")
	for i, profile := range profiles {
		fmt.Printf("%d. %s\n", i+1, profile)
	}

	fmt.Print("Select profile number: ")
	var choice int
	_, err = fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(profiles) {
		fmt.Println("Invalid choice")
		return
	}

	selectedProfile := profiles[choice-1]
	qr.CreateQR(selectedProfile)
}
