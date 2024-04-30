package model

import (
	orm "application/pkg/mysql"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Alignment struct {
	gorm.Model
	RoadName      string `json:"road_name" gorm:"type:varchar(256);comment:道路名称"`
	RoadType      string `json:"road_type" gorm:"type:varchar(256);comment:道路类型"`
	Length        int    `json:"length" gorm:"type:int;comment:长度"`
	Width         int    `json:"width" gorm:"type:int;comment:宽度"`
	RoadDetail    string `json:"road_detail" gorm:"type:varchar(256);comment:道路详情"`
	RoadStatus    string `json:"road_status" gorm:"type:varchar(256);comment:道路状态"`
	BeLongitudes  string `json:"be_longitudes" gorm:"type:varchar(256);comment:起点经度"`
	BeLatitudes   string `json:"be_latitudes" gorm:"type:varchar(256);comment:起点纬度"`
	EndLongitudes string `json:"end_longitudes" gorm:"type:varchar(256);comment:终点经度"`
	EndLatitudes  string `json:"end_latitudes" gorm:"type:varchar(256);comment:终点纬度"`
}

func (Alignment) TableName() string {
	return "alignment"
}

type AlignmentManager struct {
	orm *orm.Client
}

// CreateCollectionTasksTable 创建数据表

func CreateAlignmentTable(ctx context.Context, lar *AlignmentManager) error {
	tname := new(Alignment).TableName()
	if lar.orm.DB().Migrator().HasTable(tname) {
		fmt.Println("table already exists")
		return nil
	}
	// 迁移模型
	return lar.orm.DB().AutoMigrate(new(Alignment))
}

// CreateTable 创建数据表
func (lar *AlignmentManager) CreateTables(ctx context.Context) (err error) {
	return CreateAlignmentTable(ctx, lar)
}

type AlignmentManagerWithOption func(*AlignmentManager)

func NewAlignmentManager(orm *orm.Client, ops ...AlignmentManagerWithOption) (*AlignmentManager, error) {
	alignmentorm := &AlignmentManager{
		orm: orm,
	}
	for _, op := range ops {
		op(alignmentorm)
	}
	if err := alignmentorm.CreateTables(context.TODO()); err != nil {
		return nil, err
	}
	return alignmentorm, nil
}

func (alignmentorm *AlignmentManager) QueryCount(ctx context.Context, query *WhereQuery) (int64, error) {
	var count int64
	if err := alignmentorm.orm.DB().Model(&Alignment{}).Where(query.Query, query.Args...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (alignmentorm *AlignmentManager) QueryList(ctx context.Context, query *WhereQuery,
	ops ...QueryWithOption) (list []*Alignment, err error) {
	list = []*Alignment{}
	db := alignmentorm.orm.DB().Where(query.Query, query.Args...).Order("created_at DESC")
	for _, op := range ops {
		db = op(db)
	}
	err = db.Find(&list).Error
	return
}

// 注册用户
func (alignmentorm *AlignmentManager) Insert(ctx context.Context, item *Alignment) error {
	return alignmentorm.orm.DB().Create(item).Error
}

// 更新用户
func (alignmentorm *AlignmentManager) Update(ctx context.Context, ID uint, item *Alignment) error {
	return alignmentorm.orm.DB().Model(&Alignment{}).Where("id = ?", ID).Updates(item).Error
}

// 删除用户
func (alignmentorm *AlignmentManager) Delete(ctx context.Context, ID uint) error {
	return alignmentorm.orm.DB().Where("id = ?", ID).Delete(&Alignment{}).Error
}
