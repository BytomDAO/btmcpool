package logger

import (
	"testing"
)

func TestSimpleLog(t *testing.T) {
	err := Init(DebugLevel)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		Info("Hello univ!", "aa", 5, "zxc", 111)
	}
}

func TestComplexLog(t *testing.T) {
	err := InitWithFields(InfoLevel, map[string]interface{}{
		"a":     1,
		"b":     2,
		"aaaaa": "ccccc",
	})
	if err != nil {
		t.Error(err)
	}

	Info("Hello!", "ddd", 3)
}
