package collections

import "testing"

func TestHello(t *testing.T) {
	if got := hello(); got != "hello" {
		t.Errorf("got %q, wanted %q", got, "hello")
	}
}
