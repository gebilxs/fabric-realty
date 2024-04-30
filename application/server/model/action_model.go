package model

import (
	orm "application/pkg/mysql"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Action struct {
	gorm.Model
	ActionName   string `json:"actionName" gorm:"type:varchar(256);comment:动作名"`
	ActionType   string `json:"actionType" gorm:"type:varchar(256);comment:动作类型"`
	Driver       string `json:"driver" gorm:"type:varchar(256);comment:驾驶员"`
	Acceleration string `json:"acceleration" gorm:"type:varchar(256);comment:加速度"`
	Steering     string `json:"steering" gorm:"type:varchar(256);comment:转向"`
	Result       string `json:"result" gorm:"type:varchar(256);comment:结果"`
}

func (Action) TableName() string {
	return "action"
}

type ActionManager struct {
	orm *orm.Client
}

// CreateMusicTasksTable 创建数据表

func CreateActionTable(ctx context.Context, lar *ActionManager) error {
	tname := new(Action).TableName()
	if lar.orm.DB().Migrator().HasTable(tname) {
		fmt.Println("table already exists")
		return nil
	}
	// 迁移模型
	return lar.orm.DB().AutoMigrate(new(Action))
}

// CreateTable 创建数据表
func (lar *ActionManager) CreateTables(ctx context.Context) (err error) {
	return CreateActionTable(ctx, lar)
}

type ActionManagerWithOption func(*ActionManager)

func NewActionManager(orm *orm.Client, ops ...ActionManagerWithOption) (*ActionManager, error) {
	actionorm := &ActionManager{
		orm: orm,
	}
	for _, op := range ops {
		op(actionorm)
	}
	if err := actionorm.CreateTables(context.TODO()); err != nil {
		return nil, err
	}
	return actionorm, nil
}

func (actionorm *ActionManager) QueryCount(ctx context.Context, query *WhereQuery) (int64, error) {
	var count int64
	if err := actionorm.orm.DB().Model(&Action{}).Where(query.Query, query.Args...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (actionorm *ActionManager) QueryList(ctx context.Context, query *WhereQuery,
	ops ...QueryWithOption) (list []*Action, err error) {
	list = []*Action{}
	db := actionorm.orm.DB().Where(query.Query, query.Args...).Order("created_at DESC")
	for _, op := range ops {
		db = op(db)
	}
	err = db.Find(&list).Error
	return
}

// 注册用户
func (actionorm *ActionManager) Insert(ctx context.Context, item *Action) error {
	return actionorm.orm.DB().Create(item).Error
}

// 更新用户
func (actionorm *ActionManager) Update(ctx context.Context, ID uint, item *Action) error {
	return actionorm.orm.DB().Model(&Action{}).Where("id = ?", ID).Updates(item).Error
}

// 删除用户
func (actionorm *ActionManager) Delete(ctx context.Context, ID uint) error {
	return actionorm.orm.DB().Where("id = ?", ID).Delete(&Action{}).Error
}
