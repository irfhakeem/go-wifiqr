package cmd

import (
	"flag"
	"fmt"
	"go-wifi-qr/internal/qr"
	"go-wifi-qr/internal/utils"
	"go-wifi-qr/internal/wifi"
	"os"
)

func ParseCommand() {
	allFlag := flag.Bool("a", false, "Generate QR for all saved WiFi profiles")
	findFlag := flag.String("f", "", "Generate QR for specific SSID")
	currentFlag := flag.Bool("c", false, "Generate QR for currently connected WiFi")
	testFlag := flag.Bool("t", false, "Speedtest for current connected WiFi")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *helpFlag || len(os.Args) == 1 {
		fmt.Println("Usage:")
		fmt.Println("  wifiqr -a           Generate QR for all saved WiFi profiles")
		fmt.Println("  wifiqr -f <SSID>    Generate QR for specific SSID")
		fmt.Println("  wifiqr -c           Generate QR for current connected WiFi")
		fmt.Println("  wifiqr -t           Speedtest for current connected WiFi")
		fmt.Println("  wifiqr -h           Show help")
		return
	}

	switch {
	case *allFlag:
		wifi.ListAndSelectProfile()
	case *findFlag != "":
		fmt.Println(*findFlag)
		qr.CreateQR(*findFlag)
	case *currentFlag:
		ssid, err := wifi.GetCurrentWiFiProfile()
		if err != nil {
			fmt.Println("Error getting current WiFi profile:", err)
			return
		}
		qr.CreateQR(ssid)
	case *testFlag:
		utils.TestConnection()
	default:
		fmt.Println("Invalid option. Use -h for help.")
	}
}
