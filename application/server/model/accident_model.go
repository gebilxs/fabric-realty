package model

import (
	orm "application/pkg/mysql"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Accident struct {
	gorm.Model
	AccidentName     string `json:"accidentName" gorm:"type:varchar(256);comment:事故名"`
	AccidentType     string `json:"accidentType" gorm:"type:varchar(256);comment:事故类型"`
	Driver           string `json:"driver" gorm:"type:varchar(256);comment:驾驶员"`
	AccidentLocation string `json:"accidentLocation" gorm:"type:varchar(256);comment:事故地点"`
	AccidentNotice   string `json:"accidentNotice" gorm:"type:varchar(256);comment:事故通知"`
	AccidentResult   string `json:"accidentResult" gorm:"type:varchar(256);comment:事故结果"`
	AccidentReason   string `json:"accidentReason" gorm:"type:varchar(256);comment:事故原因"`
	AccidentStatus   string `json:"accidentStatus" gorm:"type:varchar(256);comment:事故状态"`
}

func (Accident) TableName() string {
	return "accident"
}

type AccidentManager struct {
	orm *orm.Client
}

// CreateMusicTasksTable 创建数据表

func CreateAccidentTable(ctx context.Context, lar *AccidentManager) error {
	tname := new(Accident).TableName()
	if lar.orm.DB().Migrator().HasTable(tname) {
		fmt.Println("table already exists")
		return nil
	}
	// 迁移模型
	return lar.orm.DB().AutoMigrate(new(Accident))
}

// CreateTable 创建数据表
func (lar *AccidentManager) CreateTables(ctx context.Context) (err error) {
	return CreateAccidentTable(ctx, lar)
}

type AccidentManagerWithOption func(*AccidentManager)

func NewAccidentManager(orm *orm.Client, ops ...AccidentManagerWithOption) (*AccidentManager, error) {
	accidentorm := &AccidentManager{
		orm: orm,
	}
	for _, op := range ops {
		op(accidentorm)
	}
	if err := accidentorm.CreateTables(context.TODO()); err != nil {
		return nil, err
	}
	return accidentorm, nil
}
