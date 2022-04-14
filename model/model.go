package model

import (
	"IMConnection/conf"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

var DB *gorm.DB

func Setup() {
	var err error

	// pass conf to dsn, meet the problem that there is not
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.DatabaseSetting.User,
		conf.DatabaseSetting.Password,
		conf.DatabaseSetting.Host,
		conf.DatabaseSetting.Name,
	)

	// open the database and buffer the conf
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 迁移时禁用外键
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.DatabaseSetting.TablePrefix, // set the prefix name of table
			SingularTable: true,                             // use singular table by default
		},
		Logger: logger.Default.LogMode(logger.Info), // set log mode
	})

	// some init set of database
	mysqlDB, err := DB.DB()
	if err != nil {
		log.Panicln("db.DB() err: ", err)
	}
	mysqlDB.SetMaxIdleConns(conf.DatabaseSetting.SetMaxIdleConns)       // SetMaxIdleConns 设置空闲连接池中连接的最大数量
	mysqlDB.SetMaxOpenConns(conf.DatabaseSetting.SetMaxOpenConns)       // SetMaxOpenConns 设置打开数据库连接的最大数量
	mysqlDB.SetConnMaxLifetime(conf.DatabaseSetting.SetConnMaxLifetime) // SetConnMaxLifetime 设置了连接可复用的最大时间

	// set auto migrate
	migration()
}

func migration() {
	//自动迁移模式
	DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&User{},
			&Relation{},
		)
}
