package model

import (
	orm "application/pkg/mysql"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Car struct {
	gorm.Model
	CarName     string `json:"carName" gorm:"type:varchar(256);comment:车名"`
	CarType     string `json:"carType" gorm:"type:varchar(256);comment:车类型"`
	CarColor    string `json:"carColor" gorm:"type:varchar(256);comment:车颜色"`
	CarEngine   string `json:"carEngine" gorm:"type:varchar(256);comment:车引擎"`
	CarLocation string `json:"carLocation" gorm:"type:varchar(256);comment:车位置"`
	CarHealth   string `json:"carHealth" gorm:"type:varchar(256);comment:车健康"`
}

func (Car) TableName() string {
	return "car"
}

type CarManager struct {
	orm *orm.Client
}

// CreateMusicTasksTable 创建数据表

func CreateCarTable(ctx context.Context, lar *CarManager) error {
	tname := new(Car).TableName()
	if lar.orm.DB().Migrator().HasTable(tname) {
		fmt.Println("table already exists")
		return nil
	}
	// 迁移模型
	return lar.orm.DB().AutoMigrate(new(Car))
}

// CreateTable 创建数据表
func (lar *CarManager) CreateTables(ctx context.Context) (err error) {
	return CreateCarTable(ctx, lar)
}

type CarManagerWithOption func(*CarManager)

func NewCarManager(orm *orm.Client, ops ...CarManagerWithOption) (*CarManager, error) {
	carorm := &CarManager{
		orm: orm,
	}
	for _, op := range ops {
		op(carorm)
	}
	if err := carorm.CreateTables(context.TODO()); err != nil {
		return nil, err
	}
	return carorm, nil
}
