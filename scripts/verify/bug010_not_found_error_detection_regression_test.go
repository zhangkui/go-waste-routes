package verify

import (
	"errors"
	"testing"

	"go-waste-routes/internal/service"
)

func TestBug010_NotFoundErrorIsDetected(t *testing.T) {
	if service.ErrNotFound(errors.New("record not found")) == false {
		t.Fatal("not found error was not detected")
	}
}
