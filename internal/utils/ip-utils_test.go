package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Format_Function_Must_Remove_Port_Number(t *testing.T) {
	assert.Equal(t, "192.168.0.1", FormattIp("192.168.0.1:3000"))
}

func Test_Format_Function_Should_Not_Modify_String_If_No_Port_Is_Available(t *testing.T) {
	assert.Equal(t, "192.168.0.1", FormattIp("192.168.0.1"))
}

func Test_Format_Function_Should_Remove_The_Port(t *testing.T) {
	assert.NotEqual(t, "192.168.0.1:3000", FormattIp("192.168.0.1:3000"))
}
