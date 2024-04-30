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

func (accidentorm *AccidentManager) QueryCount(ctx context.Context, query *WhereQuery) (int64, error) {
	var count int64
	if err := accidentorm.orm.DB().Model(&Accident{}).Where(query.Query, query.Args...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (accidentorm *AccidentManager) QueryList(ctx context.Context, query *WhereQuery,
	ops ...QueryWithOption) (list []*Accident, err error) {
	list = []*Accident{}
	db := accidentorm.orm.DB().Where(query.Query, query.Args...).Order("created_at DESC")
	for _, op := range ops {
		db = op(db)
	}
	err = db.Find(&list).Error
	return
}

// 注册用户
func (accidentorm *AccidentManager) Insert(ctx context.Context, item *Accident) error {
	return accidentorm.orm.DB().Create(item).Error
}

// 更新用户
func (accidentorm *AccidentManager) Update(ctx context.Context, ID uint, item *Accident) error {
	return accidentorm.orm.DB().Model(&Accident{}).Where("id = ?", ID).Updates(item).Error
}

// 删除用户
func (accidentorm *AccidentManager) Delete(ctx context.Context, ID uint) error {
	return accidentorm.orm.DB().Where("id = ?", ID).Delete(&Accident{}).Error
}
