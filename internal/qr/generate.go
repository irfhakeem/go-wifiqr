package qr

import (
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
)

func CreateQR(p string) {
	data := GetProfileDetail(p)
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
