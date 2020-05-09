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

func Test_Utils_StringToDate(t *testing.T) {
	utils := NewUtils()
	date, err := utils.StringToDate("2006-01-02T15:04:05Z")
	assert.NoError(t, err)

	time := date.Unix()
	assert.Equal(t, int64(1136214245), time)
}

func Test_Utils_StringToDate_wrongDate(t *testing.T) {
	utils := NewUtils()
	_, err := utils.StringToDate("not a date")
	assert.EqualError(t, err, "cannot parse date not a date: parsing time \"not a date\" as \"2006-01-02T15:04:05Z\": cannot parse \"not a date\" as \"2006\"")
}
