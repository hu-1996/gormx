package gormx

import (
	"errors"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(db *gorm.DB) {
	DB = db
}

func convert[T any, R any](entity *T) (*R, error) {
	convertInterface, ok := any(entity).(ConvertInterface)
	if !ok {
		return nil, errors.New("The source data is not implemented [ConvertInterface]")
	}
	return convertInterface.Convert().(*R), nil
}

func SelectById[T any](id any) (*T, error) {
	var entity T
	err := DB.Last(&entity, "id = ?", id).Error
	return &entity, err
}

func SelectConvertById[T any, R any](id any) (*R, error) {
	var entity T
	err := DB.Last(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return convert[T, R](&entity)
}

func SelectByIds[T any](ids any) ([]*T, error) {
	var records []*T
	err := DB.Find(&records, "id IN ?", ids).Error
	return records, err
}

func SelectConvertByIds[T any, R any](ids any) ([]*R, error) {
	var records []*T
	err := DB.Find(&records, "id IN ?", ids).Error
	if err != nil {
		return nil, err
	}
	var convertRecords []*R
	for _, record := range records {
		convertRecord, err := convert[T, R](record)
		if err != nil {
			return nil, err
		}
		convertRecords = append(convertRecords, convertRecord)
	}
	return convertRecords, err
}

func SelectOne[T any](query interface{}, args ...interface{}) (*T, error) {
	var entity T
	err := DB.Where(query, args...).Last(&entity).Error
	return &entity, err
}

func SelectOneConvert[T any, R any](query interface{}, args ...interface{}) (*R, error) {
	var entity T
	err := DB.Where(query, args...).Last(&entity).Error
	if err != nil {
		return nil, err
	}
	return convert[T, R](&entity)
}

func SelectList[T any](order interface{}, query interface{}, args ...interface{}) ([]*T, error) {
	var records []*T
	err := DB.Where(query, args...).Order(order).Find(&records).Error
	return records, err
}

func DBSelectList[T any](db *gorm.DB, order interface{}, query interface{}, args ...interface{}) ([]*T, error) {
	var records []*T
	err := db.Where(query, args...).Order(order).Find(&records).Error
	return records, err
}

func SelectListConvert[T any, R any](order interface{}, query interface{}, args ...interface{}) ([]*R, error) {
	var records []*T
	err := DB.Where(query, args...).Order(order).Find(&records).Error
	var convertRecords []*R
	for _, record := range records {
		convertRecord, err := convert[T, R](record)
		if err != nil {
			return nil, err
		}
		convertRecords = append(convertRecords, convertRecord)
	}
	return convertRecords, err
}

func DBSelectListConvert[T any, R any](db *gorm.DB, order interface{}, query interface{}, args ...interface{}) ([]*R, error) {
	var records []*T
	err := db.Where(query, args...).Order(order).Find(&records).Error
	var convertRecords []*R
	for _, record := range records {
		convertRecord, err := convert[T, R](record)
		if err != nil {
			return nil, err
		}
		convertRecords = append(convertRecords, convertRecord)
	}
	return convertRecords, err
}

func SelectPage[T any](current, pageSize int32, order interface{}, query interface{}, args ...interface{}) ([]*T, int64, error) {
	var records []*T
	var total int64
	err := DB.Scopes(Paginate(current, pageSize)).Where(query, args...).Order(order).Find(&records).Offset(-1).Limit(-1).Count(&total).Error
	return records, total, err
}

func DBSelectPage[T any](db *gorm.DB, current, pageSize int32, order interface{}, query interface{}, args ...interface{}) ([]*T, int64, error) {
	var records []*T
	var total int64
	err := db.Scopes(Paginate(current, pageSize)).Where(query, args...).Order(order).Find(&records).Offset(-1).Limit(-1).Count(&total).Error
	return records, total, err
}

func SelectPageConvert[T any, R any](current, pageSize int32, order interface{}, query interface{}, args ...interface{}) ([]*R, int64, error) {
	var records []*T
	var total int64
	err := DB.Scopes(Paginate(current, pageSize)).Where(query, args...).Order(order).Find(&records).Offset(-1).Limit(-1).Count(&total).Error

	var convertRecords []*R
	for _, record := range records {
		convertRecord, err := convert[T, R](record)
		if err != nil {
			return nil, 0, err
		}
		convertRecords = append(convertRecords, convertRecord)
	}
	return convertRecords, total, err
}

func DBSelectPageConvert[T any, R any](db *gorm.DB, current, pageSize int32, order interface{}, query interface{}, args ...interface{}) ([]*R, int64, error) {
	var records []*T
	var total int64
	err := db.Scopes(Paginate(current, pageSize)).Where(query, args...).Order(order).Find(&records).Offset(-1).Limit(-1).Count(&total).Error

	var convertRecords []*R
	for _, record := range records {
		convertRecord, err := convert[T, R](record)
		if err != nil {
			return nil, 0, err
		}
		convertRecords = append(convertRecords, convertRecord)
	}
	return convertRecords, total, err
}

func Count[T any](query interface{}, args ...interface{}) (int64, error) {
	var entity T
	var total int64
	err := DB.Model(&entity).Where(query, args...).Count(&total).Error
	return total, err
}

func DBCount[T any](db *gorm.DB, query interface{}, args ...interface{}) (int64, error) {
	var entity T
	var total int64
	err := db.Model(&entity).Where(query, args...).Count(&total).Error
	return total, err
}

func Exist[T any](query interface{}, args ...interface{}) (bool, error) {
	var entity T
	var total int64
	err := DB.Model(&entity).Where(query, args...).Count(&total).Error
	return total > 0, err
}

func DBExist[T any](db *gorm.DB, query interface{}, args ...interface{}) (bool, error) {
	var entity T
	var total int64
	err := db.Model(&entity).Where(query, args...).Count(&total).Error
	return total > 0, err
}

func Insert[T any](entity *T) (int64, error) {
	result := DB.Create(entity)
	return result.RowsAffected, result.Error
}

func InsertBatches[T any](entity []*T) (int64, error) {
	result := DB.Create(entity)
	return result.RowsAffected, result.Error
}

func TxInsert[T any](tx *gorm.DB, entity *T) (int64, error) {
	result := tx.Create(entity)
	return result.RowsAffected, result.Error
}

func TxInsertBatches[T any](tx *gorm.DB, entity []*T) (int64, error) {
	result := tx.Create(entity)
	return result.RowsAffected, result.Error
}

func Update[T any](entity *T) (int64, error) {
	result := DB.Save(entity)
	return result.RowsAffected, result.Error
}

func UpdateBatches[T any](entity []*T) (int64, error) {
	result := DB.Save(entity)
	return result.RowsAffected, result.Error
}

func Updates[T any](entity *T, query interface{}, args ...interface{}) (int64, error) {
	result := DB.Model(entity).Where(query, args...).Updates(entity)
	return result.RowsAffected, result.Error
}

func UpdatesMap[T any](data map[string]interface{}, query interface{}, args ...interface{}) (int64, error) {
	var entity T
	result := DB.Model(&entity).Where(query, args...).Updates(data)
	return result.RowsAffected, result.Error
}

func TxUpdate[T any](tx *gorm.DB, entity *T) (int64, error) {
	result := tx.Save(entity)
	return result.RowsAffected, result.Error
}

func TxUpdateBatches[T any](tx *gorm.DB, entity []*T) (int64, error) {
	result := tx.Save(entity)
	return result.RowsAffected, result.Error
}

func TxUpdates[T any](tx *gorm.DB, entity *T, query interface{}, args ...interface{}) (int64, error) {
	result := tx.Model(entity).Where(query, args...).Updates(entity)
	return result.RowsAffected, result.Error
}

func TxUpdatesMap[T any](tx *gorm.DB, data map[string]interface{}, query interface{}, args ...interface{}) (int64, error) {
	var entity T
	result := tx.Model(&entity).Where(query, args...).Updates(data)
	return result.RowsAffected, result.Error
}

func DeleteById[T any](id any) (int64, error) {
	var entity T
	result := DB.Delete(&entity, "id = ?", id)
	return result.RowsAffected, result.Error
}

func DeleteByIds[T any](ids any) (int64, error) {
	var entity T
	result := DB.Delete(&entity, "id IN ?", ids)
	return result.RowsAffected, result.Error
}

func Delete[T any](query interface{}, args ...interface{}) (int64, error) {
	var entity T
	result := DB.Where(query, args...).Delete(&entity)
	return result.RowsAffected, result.Error
}

func TxDeleteById[T any](tx *gorm.DB, id any) (int64, error) {
	var entity T
	result := tx.Delete(&entity, "id = ?", id)
	return result.RowsAffected, result.Error
}

func TxDeleteByIds[T any](tx *gorm.DB, ids any) (int64, error) {
	var entity T
	result := tx.Delete(&entity, "id IN ?", ids)
	return result.RowsAffected, result.Error
}

func TxDelete[T any](tx *gorm.DB, query interface{}, args ...interface{}) (int64, error) {
	var entity T
	result := tx.Where(query, args...).Delete(&entity)
	return result.RowsAffected, result.Error
}
