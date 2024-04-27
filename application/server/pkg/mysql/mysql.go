package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MysqlConfig  定义数据库连接配置
type MysqlConfig struct {
	Host     string        `json:"host"` // ip:port
	Username string        `json:"username"`
	Passwd   string        `json:"passwd"`
	DB       string        `json:"db"`
	Timeout  time.Duration `json:"timeout"` // 设置了连接可复用的最大时间
	MaxIdle  int           `json:"maxIdle"` // 设置空闲连接池中连接的最大数量
	MaxOpen  int           `json:"MaxOpen"` // 设置打开数据库连接的最大数量
	IsDebug  bool          `json:"isDebug"`
}

type Client struct {
	db   *gorm.DB
	conn *sql.DB
}

// NewOrm 创建数据库连接对象
func NewOrm(conn *MysqlConfig, orm *gorm.Config) (*Client, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", conn.Username, conn.Passwd, conn.Host, conn.DB)
	db, err := gorm.Open(mysql.Open(dsn), orm)
	if err != nil {
		return nil, err
	}
	if conn.IsDebug {
		db = db.Debug()
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if conn.Timeout == 0 {
		conn.Timeout = 15 * time.Minute
	}
	if conn.MaxIdle == 0 {
		conn.MaxIdle = 10
	}
	if conn.MaxOpen == 0 {
		conn.MaxOpen = 100
	}
	// 设置连接池
	sqlDB.SetMaxIdleConns(conn.MaxIdle)
	sqlDB.SetMaxOpenConns(conn.MaxOpen)
	sqlDB.SetConnMaxLifetime(conn.Timeout)

	return &Client{
		db:   db,
		conn: sqlDB,
	}, nil
}

func (cli *Client) DB() *gorm.DB {
	return cli.db
}

// Conn 获取通用数据库对象 sql.DB
func (cli *Client) Conn() *sql.DB {
	return cli.conn
}

func (cli *Client) Ping() error {
	return cli.conn.Ping()
}

// Stats 返回数据库统计信息
func (cli *Client) Stats() sql.DBStats {
	return cli.conn.Stats()
}

func GetGormConfig() *gorm.Config {
	return new(gorm.Config)
}
