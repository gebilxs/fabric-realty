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
