package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/unsurper/dlt645go/protocol"
	"reflect"
	"testing"
)

func TestDlt645_0x33333433(t *testing.T) {
	data := []byte{0x68, 0x76, 0x02, 0x02, 0x47, 0x02, 0x02, 0x68, 0x91, 0x08, 0x33, 0x33, 0x34, 0x33, 0x35, 0x33, 0x33, 0x33, 0xC9, 0x16}
	var message protocol.Dlt_0x33333433
	_, err := message.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual("020247020276", message.DeviceID))
	assert.True(t, reflect.DeepEqual(0.02, message.WP))
}
