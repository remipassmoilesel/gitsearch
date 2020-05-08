package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Utils_ContainsString_True(t *testing.T) {
	utils := NewUtils()
	slice := []string{"a", "b", "c"}
	assert.True(t, utils.ContainsString(slice, "c"))
}

func Test_Utils_ContainsString_False(t *testing.T) {
	utils := NewUtils()
	slice := []string{"a", "b", "c"}
	assert.False(t, utils.ContainsString(slice, "d"))
}
