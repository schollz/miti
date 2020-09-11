package click

import (
	"testing"
	"time"
)

func TestClick(t *testing.T) {
	Play(60)
	time.Sleep(3 * time.Second)
	Stop()
	time.Sleep(1 * time.Second)
}
