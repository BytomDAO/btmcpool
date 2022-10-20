package service

import (
	"testing"

	"github.com/bytom/btmcpool/common/logger"
)

func TestService(t *testing.T) {
	t.Skip("Skipping testing in CI environment")

	c := NewConfig("prod").SetLogLevel(logger.DebugLevel)
	s := New("test", c)

	s.Run(":8081")
}
