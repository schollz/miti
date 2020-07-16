package music

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCluster(t *testing.T) {
	cluster := "CEG"
	notes, err := ParseCluster(cluster)
	assert.Nil(t, err)
	assert.Equal(t, []Note{NewNote("C", 4), NewNote("E", 4), NewNote("G", 4)}, notes)
}
