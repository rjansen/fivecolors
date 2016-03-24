package data

import (
    "log"
    "fmt"
    "errors"
    "database/sql"
    //Load Mysql Driver
    _ "github.com/go-sql-driver/mysql"
    "farm.e-pedion.com/repo/fivecolors/config"
)

//Pool is a variable to hold the Database Pool
var ( 
    pool *DBPool
)

//GetPool gets the singleton db pool reference.
//You must call Setup before get the pool reference 
func GetPool() (*DBPool, error) {
    if pool == nil {
        return nil, errors.New("SetupMustCalled: Message='You must call Setup with a DBConfig before get a DBpool reference')")
    }
    return pool, nil
}

//Setup configures a poll for database connections
func Setup(config *config.DBConfig) {
    datasource := Datasource{
        Driver: config.Driver,
        Username: config.Username,
        Password: config.Password,
        URL: config.URL,
    }
    pool = &DBPool{
        MinCons: 5,
        MaxCons: 10,
        Datasource: datasource,
    }
    log.Printf("data.Setted: Config=%+v", config)
}

//DBPool controls how new sql.DB will create and maintained
type DBPool struct {
    MinCons int
    MaxCons int
    Datasource
}

//GetConnection creates and returns a new sql.DB
func (d *DBPool) GetConnection() (*sql.DB, error) {
    if d == nil {
        return nil, errors.New("SetupMustCalled: Message='You must call Setup with a DBConfig before get a DBpool reference')")
    }
    log.Printf("data.GetConnection: DBPool=%+v", d)
    conn, err := sql.Open(d.Datasource.Driver, d.Datasource.GetDSN());
    if err != nil {
        return nil, errors.New("GetConnectionError: Cause=" + err.Error())
    }
    return conn, nil 
}

//Datasource holds parameterts to create new sql.DB connections 
type Datasource struct {
    Driver string
    URL string
    Username string
    Password string
}

//GetDSN retuns a DNS representation of Datasource struct
//DSN format: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func (d *Datasource) GetDSN() string {
    return fmt.Sprintf("%s:%s@%s", d.Username, d.Password, d.URL)
}

//FromDSN fills the connection parameters of this Datasource instance 
func (d *Datasource) FromDSN(DSN string) error {
    regex := "(()?(:())@)?()?/()?"
    return errors.New("NotImplemented: Regex=" + regex)
}
