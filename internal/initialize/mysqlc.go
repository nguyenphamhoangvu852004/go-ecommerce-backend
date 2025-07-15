package initialize

import (
	// "chapapp-backend-api/internal/entity"
	// "chapapp-backend-api/internal/utils"
	// "errors"
	"database/sql"
	"fmt"
	"go-ecommerce-backend-api/global"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func checkErrC(err error, msg string) {
	if err != nil {
		global.Logger.Error(msg)
		panic(err)
	}
}

func InitMysqlC() {
	global.Logger.Info("Init mysql...")
	m := global.Config.Mysql
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, m.Username, m.Password, m.Host, m.Port, m.Dbname)
	db, err := sql.Open("mysql", s)
	checkErrC(err, "Init mysql failed")
	global.Mdbc = db
	global.Logger.Info("MysqlPool Initialize Successfully")

	// setPool
	SetPoolC()
}

func SetPoolC() {
	m := global.Config.Mysql
	sqlDB := global.Mdbc

	sqlDB.SetConnMaxIdleTime(time.Duration(m.MaxIdleConns))
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(m.ConnMaxLifeTime))
}
