package user

import (
	"time"

	"gorm.io/gorm"

	servicedb "github.com/zero-look/pkg/service/db"
)

// User represents the rpc table.
type User struct {
	ID         int64     `gorm:"column:id;primaryKey"` // 主键
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
	UserID     int64     `gorm:"column:user_id"`    // 用户 ID
	StaffName  string    `gorm:"column:staff_name"` // 用户名称
	Email      string    `gorm:"column:email"`      // 邮箱
	Password   string    `gorm:"column:password"`   // 密码
}

func (*User) TableName() string {
	return "user"
}

func (u *User) BeforeSave(db *gorm.DB) error {
	if !db.Statement.Changed("create_time") {
		u.CreateTime = time.Now()
	}
	u.UpdateTime = time.Now()
	return nil
}

func FindOneByID(db *gorm.DB, id int64) (*User, error) {
	var user User
	err := db.Table(user.TableName()).Take(&user, "id = ?", id).Error
	if err != nil {
		if servicedb.IsErrNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func FindOneForLogin(db *gorm.DB, staffName, password string) (*User, error) {
	var user User
	err := db.Table(user.TableName()).Take(&user, "staff_name = ? AND password = ?", staffName, password).Error
	if err != nil {
		if servicedb.IsErrNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
