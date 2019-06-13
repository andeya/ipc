package ipc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgp(t *testing.T) {
	var expected = Msgp{
		Mtype: 9,
		Mtext: []byte("henrylee2cn"),
	}
	ptr, textSize := expected.marshal()
	assert.Equal(t, len(expected.Mtext), textSize)
	var actual Msgp
	err := actual.unmarshal(textSize, ptr)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
