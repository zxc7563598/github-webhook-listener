package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/zxc7563598/github-webhook-listener/internal/model"
	"github.com/zxc7563598/github-webhook-listener/pkg/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

// InitSQLiteDB 初始化 SQLite 数据库
func InitSQLiteDB() {
	var err error
	// 获取数据库文件路径
	dbPath, err := utils.GetExecutableDir()
	if err != nil {
		log.Fatalf("数据库路径获取失败: %v", err)
	}
	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("无法创建数据库目录: %v", err)
	}
	// 连接数据库
	DB, err = gorm.Open(
		&sqlite.Dialector{
			DSN:        filepath.Join(dbPath, "database.db"),
			DriverName: "sqlite",
		},
		&gorm.Config{},
	)
	if err != nil {
		log.Fatalf("无法连接 SQLite 数据库: %v", err)
	}
	DB.Exec("PRAGMA journal_mode = WAL;")
	DB.Exec("PRAGMA busy_timeout = 10000;")
	log.Println("SQLite 数据库连接成功")
	// 自动创建表
	err = DB.AutoMigrate(
		&model.HealthMonitoring{},
		&model.WebhookLog{},
	)
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	log.Println("数据库自动迁移完成")
}
