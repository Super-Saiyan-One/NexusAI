package model

import (
	"fmt"
	"nexus-ai/model/log"
	"nexus-ai/model/worker"
	"nexus-ai/mysql"
	"nexus-ai/utils"

	"time"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// 初始化数据库连接和表结构
func InitGorm() error {
	sqlDB := mysql.GetDB()
	if sqlDB == nil {
		return fmt.Errorf("failed to get MySQL connection")
	}

	gormDB, err := gorm.Open(gmysql.New(gmysql.Config{
		Conn:                      sqlDB,
		DefaultStringSize:         256,  // 设置默认字符串长度
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，而不是 `rename`
		SkipInitializeWithVersion: true, // 跳过版本检查
	}), &gorm.Config{
		Logger:               logger.Default.LogMode(logger.Silent), // 添加详细日志
		DisableAutomaticPing: true,                                  // 禁用自动 ping
	})
	if err != nil {
		return fmt.Errorf("failed to initialize GORM: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db = gormDB

	// 初始化数据库表
	if err := Setup(gormDB); err != nil {
		return fmt.Errorf("failed to setup database tables: %v", err)
	}

	return nil
}

// GetDB 获取 GORM 数据库实例
func GetDB() *gorm.DB {
	return db
}

// 创建迁移辅助函数
func migrateTable(db *gorm.DB, table interface{}) error {
	tableName := db.NamingStrategy.TableName(fmt.Sprintf("%T", table))
	migrator := db.Migrator()

	exists := migrator.HasTable(table)
	if exists {
		utils.SysInfo(fmt.Sprintf("Table %s already exists, Migrating...", tableName))
		if err := migrator.AutoMigrate(table); err != nil {
			return fmt.Errorf("failed to migrate table %s: %v", tableName, err)
		}
		utils.SysInfo(fmt.Sprintf("Successfully migrated table %s", tableName))
	} else {
		utils.SysInfo(fmt.Sprintf("Creating table %s...", tableName))
		if err := migrator.AutoMigrate(table); err != nil {
			return fmt.Errorf("failed to create table %s: %v", tableName, err)
		}
		utils.SysInfo(fmt.Sprintf("Successfully created table %s", tableName))
	}
	return nil
}

// Setup 初始化所有数据库表
func Setup(db *gorm.DB) error {
	// 先尝试创建数据库（如果不存在）
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", mysql.GetConfig().Basic.Database)
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	// 确保使用正确的数据库
	if err := db.Exec(fmt.Sprintf("USE %s", mysql.GetConfig().Basic.Database)).Error; err != nil {
		return fmt.Errorf("failed to use database: %v", err)
	}

	// 设置表选项
	db = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci")

	// 第一批：基础表（无外键依赖）
	baseTables := []interface{}{
		&UserGroup{},
		&ChannelGroup{},
		&ModelGroup{},
		&User{},                 // 用户表最基础
		&Channel{},              // 渠道表基础
		&Model{},                // 模型表基础
		&worker.WorkerCluster{}, // 集群表基础
	}

	for _, table := range baseTables {
		if err := migrateTable(db, table); err != nil {
			return fmt.Errorf("failed to migrate base table: %v", err)
		}
	}

	// 第二批：依赖基础表的表
	secondaryTables := []interface{}{
		&Token{},              // 依赖 User
		&worker.WorkerGroup{}, // 依赖 WorkerCluster
	}

	for _, table := range secondaryTables {
		if err := migrateTable(db, table); err != nil {
			return fmt.Errorf("failed to migrate secondary table: %v", err)
		}
	}

	// 第三批：依赖第二批表的表
	tertiaryTables := []interface{}{
		&worker.WorkerNode{}, // 依赖 WorkerGroup 和 WorkerCluster
		&Usage{},             // 依赖 Token, User, Channel, Model
		&Quota{},             // 依赖 User
		&Payment{},           // 依赖 User
		&MessageSave{},       // 依赖 User, Token, Model, Channel
		&Billing{},           // 依赖 User
		&TaskVideo{},         // 依赖 User, Token, Model, Channel
	}

	for _, table := range tertiaryTables {
		if err := migrateTable(db, table); err != nil {
			return fmt.Errorf("failed to migrate tertiary table: %v", err)
		}
	}

	// 第四批：日志表（依赖多个其他表）
	logTables := []interface{}{
		&log.GatewayLog{},
		&log.MasterLog{},
		&log.MasterMySQLLog{},
		&log.RedisLog{},
		&log.RedisPersistLog{},
		&log.RelayLog{},
		&log.RequestLog{},
		&log.WorkerLog{},
		&log.WorkerMySQLLog{},
	}

	for _, table := range logTables {
		if err := migrateTable(db, table); err != nil {
			return fmt.Errorf("failed to migrate log table: %v", err)
		}
	}

	return nil
}
