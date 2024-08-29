package logger

import "testing"

func TestInitialize(t *testing.T) {
	err := Initialize("WARN")
	if err != nil {
		t.Errorf("initialize logger error")
	}
}
