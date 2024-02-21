// Package servicedb Pack
package servicedb

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	defaultCharset   = "utf8mb4"
	defaultParseTime = true
	defaultLoc       = "Local"
)

// Config for database config
type Config struct {
	Host string
	Port int
	Name string
	User string
	Pass string
}

func (c Config) getDNS() (string, error) {
	if c.Host == "" || c.Port == 0 || c.Name == "" || c.User == "" {
		return "", errors.New("db config should not be empty")
	}
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User, c.Pass, c.Host, c.Port, c.Name, defaultCharset, defaultParseTime, defaultLoc)
	return dns, nil
}

// New init database connection
func New(conf *Config) (*gorm.DB, error) {
	dns, err := conf.getDNS()
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dns,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		CreateBatchSize: 200, // 批量插入限制
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// IsErrNoRows 是否是未查询到数据
func IsErrNoRows(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) //nolint:bannedfunc
}

// IsRedisNil redis key 是否存在
func IsRedisNil(err error) bool {
	return errors.Is(err, redis.Nil) //nolint:bannedfunc
}

// Exist record exist
func Exist(db *gorm.DB, condition ...any) (bool, error) {
	var obj map[string]interface{}
	var err error
	if len(condition) > 1 {
		err = db.Where(condition[0], condition[1:]).Take(&obj).Error
	} else {
		err = db.Where(condition).Take(&obj).Error
	}
	if err != nil {
		if IsErrNoRows(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
