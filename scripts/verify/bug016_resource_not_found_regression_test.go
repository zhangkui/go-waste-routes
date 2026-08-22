package verify

import (
	"errors"
	"testing"

	"go-waste-routes/internal/service"
)

func TestBug016_ResourceNotFoundClassifierRecognizesLookupMiss(t *testing.T) {
	if !service.ErrNotFound(errors.New("not found")) {
		t.Fatal("lookup miss was not classified as not found")
	}
}
