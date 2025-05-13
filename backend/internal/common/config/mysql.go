package config

import (
	"fmt"
	stdlog "log"
	"time"

	"github.com/fynn404/gin-demo/backend/internal/model"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	xormlog "xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

var DB *xorm.Engine

// InitMySQL 初始化MySQL连接
func InitMySQL() {
	dbCfg := GlobalConfig.Database

	// 构建DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
	)

	// 创建数据库引擎
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		stdlog.Fatalf("Failed to create MySQL engine: %v", err)
	}

	// 设置连接池
	engine.SetMaxIdleConns(dbCfg.MaxIdleConns)                     // 设置空闲连接池中的最大连接数
	engine.SetMaxOpenConns(dbCfg.MaxOpenConns)                     // 设置打开的最大连接数
	engine.SetConnMaxLifetime(dbCfg.ConnMaxLifetime * time.Second) // 设置连接可复用的最大时间

	// 设置名称映射规则（数据库表名和字段名使用蛇形命名）
	engine.SetMapper(names.GonicMapper{})

	// 设置日志级别
	if GlobalConfig.Server.Mode == "debug" {
		engine.ShowSQL(true) // 在调试模式下打印SQL语句
		engine.Logger().SetLevel(xormlog.LOG_DEBUG)
	}

	// 同步数据库结构
	if err := syncModels(engine); err != nil {
		stdlog.Fatalf("Failed to sync database schema: %v", err)
	}

	// 测试连接
	if err := engine.Ping(); err != nil {
		stdlog.Fatalf("Failed to connect to MySQL: %v", err)
	}

	DB = engine
	stdlog.Printf("MySQL connected successfully. MaxOpenConns: %d, MaxIdleConns: %d",
		dbCfg.MaxOpenConns, dbCfg.MaxIdleConns)
}

// GetDB 获取数据库实例
func GetDB() *xorm.Engine {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			stdlog.Printf("Error closing MySQL connection: %v", err)
		}
	}
}

// syncModels 同步数据库模型
func syncModels(engine *xorm.Engine) error {
	// 开始事务
	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	// 临时禁用外键检查
	if _, err := session.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		session.Rollback()
		return err
	}

	// 同步所有模型
	if err := session.Sync2(
		new(model.User),
		new(model.TodoItem),
	); err != nil {
		session.Rollback()
		return err
	}

	// 重新启用外键检查
	if _, err := session.Exec("SET FOREIGN_KEY_CHECKS=1"); err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}

// Transaction 事务处理
func Transaction(f func(*xorm.Session) error) error {
	session := DB.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	if err := f(session); err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}
