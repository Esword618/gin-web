package tests

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	m "github.com/go-sql-driver/mysql"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/fsm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Mysql() {
	cfg, err := m.ParseDSN(global.Conf.Mysql.Uri)
	if err != nil {
		panic(fmt.Sprintf("[unit test]initialize mysql failed: %v", err))
	}
	global.Conf.Mysql.DSN = *cfg

	global.Log.Info(ctx, "[unit test]mysql dsn: %s", cfg.FormatDSN())
	var l glogger.Interface
	if global.Conf.Mysql.NoSql {
		// not show sql log
		l = global.Log.LogMode(glogger.Silent)
	} else {
		l = global.Log.LogMode(glogger.Info)
	}
	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   global.Conf.Mysql.TablePrefix + "_",
			SingularTable: true,
		},
		// select * from xxx => select a,b,c from xxx
		QueryFields: true,
		Logger:      l,
	})
	if err != nil {
		panic(fmt.Sprintf("[unit test]initialize mysql failed: %v", err))
	}
	global.Mysql = db
	autoMigrate()
	global.Log.Debug(ctx, "[unit test]initialize mysql success")
}

func autoMigrate() {
	global.Mysql.WithContext(ctx).AutoMigrate(
		new(ms.SysMenu),
		new(ms.SysMenuRoleRelation),
		new(ms.SysApi),
		new(ms.SysCasbin),
		new(ms.SysOperationLog),
		new(ms.SysMessage),
		new(ms.SysMessageLog),
		new(ms.SysMachine),
		new(ms.SysDict),
		new(ms.SysDictData),
		new(models.Leave),
		new(models.SysUser),
		new(models.SysRole),
	)
	// auto migrate fsm
	fsm.Migrate(global.Mysql, fsm.WithCtx(ctx))
}
