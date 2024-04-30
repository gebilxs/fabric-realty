package model

import (
	orm "application/pkg/mysql"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Music struct {
	gorm.Model
	MusicName   string `json:"musicName" gorm:"type:varchar(256);comment:音乐名"`
	MusicKey    string `json:"musicKey" gorm:"type:varchar(256);comment:音乐key"`
	MusicAuthor string `json:"musicAuthor" gorm:"type:varchar(256);comment:音乐作者"`
	MusicType   string `json:"musicType" gorm:"type:varchar(256);comment:音乐类型"`
	MusicAlbum  string `json:"musicAlbum" gorm:"type:varchar(256);comment:音乐专辑"`
}

func (Music) TableName() string {
	return "music"
}

type MusicManager struct {
	orm *orm.Client
}

// CreateMusicTasksTable 创建数据表

func CreateMusicTable(ctx context.Context, lar *MusicManager) error {
	tname := new(Music).TableName()
	if lar.orm.DB().Migrator().HasTable(tname) {
		fmt.Println("table already exists")
		return nil
	}
	// 迁移模型
	return lar.orm.DB().AutoMigrate(new(Music))
}

// CreateTable 创建数据表
func (lar *MusicManager) CreateTables(ctx context.Context) (err error) {
	return CreateMusicTable(ctx, lar)
}

type MusicManagerWithOption func(*MusicManager)

func NewMusicManager(orm *orm.Client, ops ...MusicManagerWithOption) (*MusicManager, error) {
	musicorm := &MusicManager{
		orm: orm,
	}
	for _, op := range ops {
		op(musicorm)
	}
	if err := musicorm.CreateTables(context.TODO()); err != nil {
		return nil, err
	}
	return musicorm, nil
}

func (musicorm *MusicManager) QueryCount(ctx context.Context, query *WhereQuery) (int64, error) {
	var count int64
	if err := musicorm.orm.DB().Model(&Music{}).Where(query.Query, query.Args...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (musicorm *MusicManager) QueryList(ctx context.Context, query *WhereQuery,
	ops ...QueryWithOption) (list []*Music, err error) {
	list = []*Music{}
	db := musicorm.orm.DB().Where(query.Query, query.Args...).Order("created_at DESC")
	for _, op := range ops {
		db = op(db)
	}
	err = db.Find(&list).Error
	return
}

// 注册用户
func (musicorm *MusicManager) Insert(ctx context.Context, item *Music) error {
	return musicorm.orm.DB().Create(item).Error
}

// 更新用户
func (musicorm *MusicManager) Update(ctx context.Context, ID uint, item *Music) error {
	return musicorm.orm.DB().Model(&Music{}).Where("id = ?", ID).Updates(item).Error
}

// 删除用户
func (musicorm *MusicManager) Delete(ctx context.Context, ID uint) error {
	return musicorm.orm.DB().Where("id = ?", ID).Delete(&Music{}).Error
}
