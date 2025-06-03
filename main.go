package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mdp/qrterminal/v3"
)

func createQR(p string) {
	data := getProfileDetail(p)
	if data["authentication"] == "" && data["password"] == "" {
		fmt.Printf("No data found for profile: %s\n", p)
		return
	}

	password := data["password"]
	auth := data["authentication"]
	payload := fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", auth, p, password)

	config := qrterminal.Config{
		HalfBlocks: true,
		Level:      qrterminal.M,
		Writer:     os.Stdout,
	}

	qrterminal.GenerateWithConfig(payload, config)
	fmt.Println("password:", password)
}

func parseProfiles(output string) []string {
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

func getProfileDetail(p string) map[string]string {
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

func getCurrentWiFiProfile() (string, error) {
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

func listAndSelectProfile() {
	cmd := exec.Command("netsh", "wlan", "show", "profiles")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting WiFi profiles: %v\n", err)
		return
	}

	profiles := parseProfiles(string(output))
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
	createQR(selectedProfile)
}

func main() {
	allFlag := flag.Bool("a", false, "Generate QR for all saved WiFi profiles")
	findFlag := flag.String("f", "", "Generate QR for specific SSID")
	currentFlag := flag.Bool("c", false, "Generate QR for currently connected WiFi")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *helpFlag || len(os.Args) == 1 {
		fmt.Println("Usage:")
		fmt.Println("  wifiqr -a           Generate QR for all saved WiFi profiles")
		fmt.Println("  wifiqr -f <SSID>    Generate QR for specific SSID")
		fmt.Println("  wifiqr -c           Generate QR for current connected WiFi")
		fmt.Println("  wifiqr -h           Show help")
		return
	}

	switch {
	case *allFlag:
		listAndSelectProfile()
	case *findFlag != "":
		createQR(*findFlag)
	case *currentFlag:
		ssid, err := getCurrentWiFiProfile()
		if err != nil {
			fmt.Println("Error getting current WiFi profile:", err)
			return
		}
		createQR(ssid)
	default:
		fmt.Println("Invalid option. Use -h for help.")
	}
}
