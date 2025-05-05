package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
	"time"
)

// NewMySQLDatabase 加载环境变量并初始化一个 *sql.DB 连接池
func NewMySQLDatabase() (*sql.DB, error) {
	// 1. 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// 2. 构造 DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 3. 打开数据库连接（底层会创建和管理连接池）
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening mysql: %w", err)
	}

	// 4. 配置连接池参数（可根据业务负载调整）
	db.SetMaxOpenConns(25)                 // 最多保持 25 个打开连接
	db.SetMaxIdleConns(25)                 // 最多保持 25 个空闲连接
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最大复用时间

	// 5. 测试连通性
	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("error pinging mysql: %w", err)
	}

	return db, nil
}
