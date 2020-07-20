package music

import (
	"fmt"
	"testing"
)

func TestParseCluster(t *testing.T) {
	cluster := "CEG"
	notes, err := ParseCluster(cluster)
	if err != nil {
		t.Errorf("err: %s", err.Error())
	}
	fmt.Println(notes)
}
