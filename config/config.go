package config

import (
	"github.com/rjansen/l"
	"github.com/rjansen/migi"
	raizelSQL "github.com/rjansen/raizel/sql"
	//"github.com/rjansen/avalon/identity"
	"fmt"
	"sync"
)

var (
	once sync.Once
	//Value is the currently state of the system configuration values
	Value *Configuration
)

//Configuration holds all possible configurations structs
type Configuration struct {
	Version     string          `mapstructure:"version"`
	Environment string          `mapstructure:"environment"`
	DB          DBConfig        `mapstructure:"db"`
	Handler     HandlerConfig   `mapstructure:"handler"`
	L           l.Configuration `mapstructure:"l"`
	// Identity    identity.Configuration  `mapstructure:"identity"`
	Raizel raizelSQL.Configuration `mapstructure:"raizel"`
}

func (c Configuration) String() string {
	//return fmt.Sprintf("Configuration[Version=%v ProxyConfig=%v DBConfig=%v SecurityConfig=%v CacheConfig=%v LoggerConfig=%v]", c.Version, c.ProxyConfig, c.DBConfig, c.SecurityConfig, c.CacheConfig, c.LoggerConfig)
	return fmt.Sprintf("Configuration Version=%s Environment=%s L=%s Handler=%s DB=%s Raizel=%s",
		c.Version, c.Environment,
		c.L.String(),
		c.Handler.String(),
		c.DB.String(),
		// c.Identity.String(),
		c.Raizel.String(),
	)
}

//Get returns the configuration struct
func Get() *Configuration {
	once.Do(func() {
		if Value == nil {
			if err := Setup(); err != nil {
				panic(err)
			}
		}
	})
	return Value
}

//Setup initializes the package
func Setup() error {
	migi.SetEnvPrefix("")
	migi.BindEnv("handler.port", "PORT")
	migi.BindEnv("raizel.sql.url", "DATABASE_URL")
	return migi.Unmarshal(&Value)
}
