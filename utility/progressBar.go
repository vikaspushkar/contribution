package utility

import (
	"fmt"
	"time"
)

//StopProgressWheel sends a signal to break the loop of progresswheel
var StopProgressWheel = make(chan bool)

//ShowProgressWheel prints the progressWheel till we get all the contributors
func ShowProgressWheel() {
	progressChars := []string{"\r|", "\r/", "\r-", "\r\\"}

	for ii := 0; ; ii++ {
		select {
		case <-StopProgressWheel:
			fmt.Print("\r\r\r\r\r")
			return
		default:
			fmt.Print("    ", progressChars[ii])
			time.Sleep(1 * time.Second)
			if ii >= 3 {
				ii = 0
			}
		}
	}
}
