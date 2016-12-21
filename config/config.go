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
	AssetDir    string          `mapstructure:"assetDir"`
	WebDir      string          `mapstructure:"webDir"`
	Handler     HandlerConfig   `mapstructure:"handler"`
	L           l.Configuration `mapstructure:"l"`
	// Identity    identity.Configuration  `mapstructure:"identity"`
	Raizel raizelSQL.Configuration `mapstructure:"raizel"`
}

func (c Configuration) String() string {
	return fmt.Sprintf("Configuration Version=%s Environment=%s AssetDir=%s L=%s Handler=%s Raizel=%s",
		c.Version, c.Environment, c.AssetDir,
		c.L.String(),
		c.Handler.String(),
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
	migi.BindEnv("raizel.url", "DATABASE_URL")
	return migi.Unmarshal(&Value)
}
