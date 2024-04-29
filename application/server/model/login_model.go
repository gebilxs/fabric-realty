package model

import (
	orm "application/pkg/mysql"
	"fmt"

	"context"

	"gorm.io/gorm"
)

type LoginAndRegister struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(256);uniqueIndex;comment:用户名"`
	Password string `json:"password" gorm:"type:varchar(256);comment:密码"`
	NickName string `json:"nickName" gorm:"type:varchar(256);comment:昵称"`
	Phone    string `json:"phone" gorm:"type:varchar(256);comment:手机号"`
	Email    string `json:"email" gorm:"type:varchar(256);comment:邮箱"`
	Sex      string `json:"sex" gorm:"type:varchar(256);comment:性别"`
	Address  string `json:"address" gorm:"type:varchar(256);comment:地址"`
	Age      int    `json:"age" gorm:"type:int;comment:年龄"`
}

func (LoginAndRegister) TableName() string {
	return "login_register"
}

type LoginAndRegisterManager struct {
	orm *orm.Client
}

// CreateCollectionTasksTable 创建数据表

func CreateLoginTable(ctx context.Context, lar *LoginAndRegisterManager) error {
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
	return CreateLoginTable(ctx, lar)
}

type LoginTaskManagerWithOption func(*LoginAndRegisterManager)

func NewLoginTaskManager(orm *orm.Client, ops ...LoginTaskManagerWithOption) (*LoginAndRegisterManager, error) {
	loginorm := &LoginAndRegisterManager{
		orm: orm,
	}
	for _, op := range ops {
		op(loginorm)
	}
	if err := loginorm.CreateTables(context.TODO()); err != nil {
		return nil, err
	}
	return loginorm, nil
}

func (loginorm *LoginAndRegisterManager) QueryCount(ctx context.Context, query *WhereQuery) (int64, error) {
	var count int64
	if err := loginorm.orm.DB().Model(&LoginAndRegister{}).Where(query.Query, query.Args...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (etm *LoginAndRegisterManager) QueryList(ctx context.Context, query *WhereQuery,
	ops ...QueryWithOption) (list []*LoginAndRegister, err error) {
	list = []*LoginAndRegister{}
	db := etm.orm.DB().Where(query.Query, query.Args...).Order("created_at DESC")
	for _, op := range ops {
		db = op(db)
	}
	err = db.Find(&list).Error
	return
}
