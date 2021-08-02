package conf

import (
	beego "github.com/beego/beego/v2/server/web"
)

// ConfigStr ...
func ConfigStr(key string) string {
	str, _ := beego.AppConfig.String(key)
	return str
}

// ConfigInt ...
func ConfigInt(key string) int {
	res, err := beego.AppConfig.Int(key)
	if err != nil {
		return 0
	}
	return res
}

// SQLConfig ...
type SQLConfig struct {
	UserName string
	Password string
	Database string
	Address  string
}

// SQLConf ...
func SQLConf() *SQLConfig {
	return &SQLConfig{
		UserName: ConfigStr("mysql_username"),
		Password: ConfigStr("mysql_password"),
		Database: ConfigStr("mysql_database"),
		Address:  ConfigStr("mysql_address"),
	}
}
