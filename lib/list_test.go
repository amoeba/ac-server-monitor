package lib

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStatusUp(t *testing.T) {
	result := getStatusMessage(true, nil)
	assert.Equal(t, result, UP)
}

func TestGetStatusDown(t *testing.T) {
	result := getStatusMessage(false, nil)
	assert.Equal(t, result, DOWN)
}

func TestGetStatusIOTimeout(t *testing.T) {
	timeoutErr := errors.New("i/o timeout")
	result := getStatusMessage(false, timeoutErr)
	assert.Equal(t, result, DOWN)
}

func TestGetStatusConnRefused(t *testing.T) {
	connRefusedErr := errors.New("read: connection refused")
	result := getStatusMessage(false, connRefusedErr)
	assert.Equal(t, result, DOWN)
}

func TestGetStatusError(t *testing.T) {
	fallbackErr := errors.New("???")
	result := getStatusMessage(false, fallbackErr)
	assert.Equal(t, result, ERROR)
}
