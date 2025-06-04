package utils

import (
	"fmt"
	"go-wifi-qr/indicator"
	"io"
	"net"
	"net/http"
	"time"
)

func TestConnection() {
	fmt.Print("Ping test ")
	done := make(chan bool)
	go indicator.ShowSpinner(done)

	start := time.Now()
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", 2*time.Second)
	done <- true

	if err != nil {
		fmt.Printf("\rPing Failed: %v\n", err)
	} else {
		latency := time.Since(start)
		fmt.Printf("\rPing: %v\n", latency)
		conn.Close()
	}

	fmt.Print("Download test ")
	done = make(chan bool)
	go indicator.ShowSpinner(done)

	url := "https://sin-speed.hetzner.com/100MB.bin"
	start = time.Now()
	resp, err := http.Get(url)
	if err != nil {
		done <- true
		fmt.Printf("\rDownload test failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	n, _ := io.Copy(io.Discard, resp.Body)
	done <- true

	duration := time.Since(start).Seconds()
	speedMbps := (float64(n) * 8) / (duration * 1e6)

	fmt.Printf("\rDownload Speed: %.2f Mbps\n", speedMbps)
}
