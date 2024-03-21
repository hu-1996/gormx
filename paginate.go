package gormx

import "gorm.io/gorm"

func Paginate(pageIndex, pageSize int32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageIndex == 0 {
			pageIndex = 1
		}
		if pageSize == 0 {
			pageSize = 10
		}
		offset := (pageIndex - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}
