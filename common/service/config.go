package service

import (
	"log"

	"github.com/bytom/btmcpool/common/logger"
)

type mode string

const (
	modeDev  mode = "dev"
	modeTest mode = "test"
	modeProd mode = "prod"
)

type Config struct {
	mode mode
	log  *logConfig
}

type logConfig struct {
	level logger.Level
}

func NewConfig(modeStr string) *Config {
	l := logger.DebugLevel
	m := mode(modeStr)

	switch m {
	case modeProd:
		l = logger.InfoLevel
	case modeDev, modeTest:
	default:
		log.Fatalf("unrecognized mode %v", m)
	}

	return &Config{
		mode: m,
		log:  &logConfig{level: l},
	}
}

func (c *Config) SetLogLevel(level logger.Level) *Config {
	c.log.level = level
	return c
}
