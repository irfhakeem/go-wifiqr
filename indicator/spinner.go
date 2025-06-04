package indicator

import (
	"fmt"
	"time"
)

func ShowSpinner(done chan bool) {
	spinner := []rune{'|', '/', '-', '\\'}
	i := 0
	for {
		select {
		case <-done:
			return
		default:
			fmt.Printf("%c", spinner[i])
			i = (i + 1) % len(spinner)
			time.Sleep(100 * time.Millisecond)
			fmt.Print("\b")
		}
	}
}
