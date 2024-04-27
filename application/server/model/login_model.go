package model

import (
	orm "application/pkg/mysql"
	"fmt"

	"context"

	"gorm.io/gorm"
)

type LoginAndRegister struct {
	gorm.Model
	username string `json:"username" gorm:"type:varchar(256);uniqueIndex;comment:用户名"`
	password string `json:"password" gorm:"type:varchar(256);comment:密码"`
	nickName string `json:"nickName" gorm:"type:varchar(256);comment:昵称"`
	phone    string `json:"phone" gorm:"type:varchar(256);comment:手机号"`
	email    string `json:"email" gorm:"type:varchar(256);comment:邮箱"`
	sex      string `json:"sex" gorm:"type:varchar(256);comment:性别"`
	address  string `json:"address" gorm:"type:varchar(256);comment:地址"`
	age      int    `json:"age" gorm:"type:int;comment:年龄"`
}

func (LoginAndRegister) TableName() string {
	return "login_register"
}

type LoginAndRegisterManager struct {
	orm *orm.Client
}

// CreateCollectionTasksTable 创建数据表

func CreateCollectionTasksTable(ctx context.Context, lar *LoginAndRegisterManager) error {
	tname := new(LoginAndRegister).TableName()
	if lar.orm.DB().Migrator().HasTable(tname) {
		fmt.Println("table already exists")
		return nil
	}
	// 迁移模型
	return lar.orm.DB().AutoMigrate(new(LoginAndRegister))
}

// CreateTable 创建数据表
func (lar *LoginAndRegisterManager) CreateTables(ctx context.Context) (err error) {
	return CreateCollectionTasksTable(ctx, lar)
}
