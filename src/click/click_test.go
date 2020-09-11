package click

import (
	"fmt"
	"testing"
	"time"
)

func TestClick(t *testing.T) {
	fmt.Println("playing")
	Play(60)
	time.Sleep(2 * time.Second)
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("click!")
		activate = true
	}
	fmt.Println("resttting")
	Reset()
	time.Sleep(2 * time.Second)
	fmt.Println("resttting")
	Reset()
	time.Sleep(2 * time.Second)
	Stop()
	time.Sleep(1 * time.Second)
}
