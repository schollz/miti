package midi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMIDI(t *testing.T) {
	assert.Nil(t, Init())
}
