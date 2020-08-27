package sad_test

import (
	"testing"

	"github.com/jswny/sad"
)

func TestHello(t *testing.T) {
	s := sad.Hello()
	expected := "Hello, world."
	if s != expected {
		t.Errorf("Expected %q, got %q instead", expected, s)
	}
}
