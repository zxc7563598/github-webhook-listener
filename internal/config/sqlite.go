package config

import (
	"log"
	"path/filepath"

	"github.com/zxc7563598/github-webhook-listener/internal/model"
	"github.com/zxc7563598/github-webhook-listener/pkg/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// InitSQLiteDB 初始化 SQLite 数据库并返回 DB 实例
func InitSQLiteDB() *gorm.DB {
	// 获取可执行文件所在目录作为数据库路径
	exeDir, err := utils.GetExecutableDir()
	if err != nil {
		log.Fatalf("数据库路径获取失败: %v", err)
	}
	dsn := filepath.Join(exeDir, "database.db")

	db, err := gorm.Open(
		&sqlite.Dialector{
			DSN:        dsn,
			DriverName: "sqlite",
		},
		&gorm.Config{},
	)
	if err != nil {
		log.Fatalf("无法连接 SQLite 数据库: %v", err)
	}
	db.Exec("PRAGMA journal_mode = WAL;")
	db.Exec("PRAGMA busy_timeout = 10000;")
	log.Println("SQLite 数据库连接成功")

	// 自动创建表
	err = db.AutoMigrate(
		&model.HealthMonitoring{},
		&model.WebhookLog{},
	)
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	log.Println("数据库自动迁移完成")
	return db
}
