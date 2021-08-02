package database

import (
	"fmt"
	"log"
	"time"

	"casbin-demo/conf"

	"github.com/astaxie/beego/orm"
)

// InitConnection ...
func InitConnection(cfg *conf.SQLConfig) error {
	// Database alias.
	dbAlias := "default"
	dbDriver := "mysql"
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.UserName, cfg.Password, cfg.Address, cfg.Database,
	)

	orm.RegisterDriver(dbDriver, orm.DRMySQL)
	orm.RegisterDataBase(dbAlias, dbDriver, connectionString)
	orm.DefaultTimeLoc = time.UTC
	orm.Debug = true

	o := orm.NewOrm()
	o.Using(dbAlias) // Using default, you can use other database

	// Drop table and re-create.
	force := false

	// Print log.
	verbose := false
	err := orm.RunSyncdb(dbAlias, force, verbose)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
