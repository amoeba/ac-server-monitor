package lib

import (
	"errors"
	"testing"
)

func TestGetStatusUp(t *testing.T) {
	var result string

	// UP
	result = getStatusMessage(true, nil)

	if result != UP {
		t.Fail()
	}
}
func TestGetStatusDown(t *testing.T) {

	// DOWN
	result := getStatusMessage(false, nil)

	if result != DOWN {
		t.Fail()
	}
}
func TestGetStatusIOTimeout(t *testing.T) {

	// I/O Timeout
	timeoutErr := errors.New("i/o timeout")

	result := getStatusMessage(false, timeoutErr)

	if result != DOWN {
		t.Fail()
	}
}
func TestGetStatusConnRefused(t *testing.T) {

	// Connection Refused
	connRefusedErr := errors.New("read: connection refused")

	result := getStatusMessage(false, connRefusedErr)

	if result != DOWN {
		t.Fail()
	}
}
func TestGetStatusError(t *testing.T) {

	// Error
	// Connection Refused
	fallbackErr := errors.New("???")

	result := getStatusMessage(false, fallbackErr)

	if result != ERROR {
		t.Fail()
	}
}
