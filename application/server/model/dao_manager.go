package model

import "gorm.io/gorm"

type WhereQuery struct {
	Query string
	Args  []any
}

const pageSize = 10

type QueryWithOption func(tx *gorm.DB) *gorm.DB

func QueryWithPage(index, size int) QueryWithOption {
	if index < 1 {
		index = 1
	}
	if size < 1 {
		size = pageSize
	}
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Limit(size).Offset((index - 1) * size)
	}
}
