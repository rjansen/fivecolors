package config

import (
	"fmt"
)

//DBConfig holds database connections parameters
type DBConfig struct {
	Driver   string
	URL      string
	Username string
	Password string
}

func (c DBConfig) String() string {
	return fmt.Sprintf("DBConfig Driver=%v URL=%v Username=%v Password=%v", c.Driver, c.URL, c.Username, c.Password)
}
