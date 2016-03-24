package config

import (
    "flag"
    "fmt"
)

var (
    handlerConfig *HandlerConfig
)

//HandlerConfig holds http handler parameters
type HandlerConfig struct {
    Version string
    BindAddress string
}

func (h *HandlerConfig) String() string {
    return fmt.Sprintf("HandlerConfig[BindAddress[%v]]", h.BindAddress)
}

//GetHandlerConfiguration gets and binds, only if necessary, parameters for http handlers
func GetHandlerConfiguration() *HandlerConfig {
    if handlerConfig == nil {
        handlerConfig = &HandlerConfig{}
        flag.StringVar(&handlerConfig.Version, "version", "0.0.1-staticversion", "Fivecolors Version")
        flag.StringVar(&handlerConfig.BindAddress, "bind_address", ":8088", "HTTP Bind Address")
    }
    return handlerConfig
}

