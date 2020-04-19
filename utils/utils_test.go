package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Utils_ContainsString_True(t *testing.T) {
	slice := []string{"a", "b", "c"}
	assert.True(t, ContainsString(slice, "c"))
}

func Test_Utils_ContainsString_False(t *testing.T) {
	slice := []string{"a", "b", "c"}
	assert.False(t, ContainsString(slice, "d"))
}
