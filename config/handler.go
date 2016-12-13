package config

import (
	"fmt"
)

//HandlerConfig holds http handler parameters
type HandlerConfig struct {
	Version string `mapstructure:"version"`
	IP      string `mapstructure:"ip"`
	Port    string `mapstructure:"port"`
}

//BindAddress returns the ip + port tcp address for socket bind purposes
func (h HandlerConfig) BindAddress() string {
	return fmt.Sprintf("%s:%s", h.IP, h.Port)
}

func (h HandlerConfig) String() string {
	return fmt.Sprintf("config.HandlerConfig Version=%s IP=%s Port=%s", h.Version, h.IP, h.Port)
}
