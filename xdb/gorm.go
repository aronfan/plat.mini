package xdb

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	NotDeleted = "deleted_at IS NULL"
)

func BuildFieldQuery[T any](ydb *gorm.DB, newObj *T, queryFields []string) (*gorm.DB, error) {
	val := reflect.ValueOf(newObj).Elem()
	typ := val.Type()

	query := ydb
	for _, fieldName := range queryFields {
		field, ok := typ.FieldByName(fieldName)
		if !ok {
			return nil, errors.New("Field " + fieldName + " not found in Struct")
		}
		tag := field.Tag.Get("gorm")
		if tag == "" {
			continue
		}
		tags := strings.Split(tag, ";")
		for _, t := range tags {
			if strings.HasPrefix(t, "column:") {
				column := strings.TrimPrefix(t, "column:")
				if len(column) > 0 {
					query = query.Where(column+" = ?", val.FieldByName(fieldName).Interface())
					break
				}
			}
		}
	}
	return query, nil
}

func Insert[T any](ydb *gorm.DB, newObj *T, queryFields []string) (*T, bool, error) {
	query, err := BuildFieldQuery(ydb, newObj, queryFields)
	if err != nil {
		return nil, false, err
	}

	err = query.First(newObj).Error
	if err == nil {
		// already exist
		return newObj, true, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// insert it
		err = ydb.Create(newObj).Error
		if err == nil {
			return newObj, false, nil
		}
		myerr, ok := err.(*mysql.MySQLError)
		if !ok {
			return nil, false, err
		}
		if myerr.Number != 1062 {
			return nil, false, err
		} else {
			// duplicated, already exist
			query, err = BuildFieldQuery(ydb, newObj, queryFields)
			if err != nil {
				return nil, false, err
			}
			err = query.First(newObj).Error
			if err != nil {
				return nil, false, err
			} else {
				return newObj, true, nil
			}
		}
	} else {
		return nil, false, err
	}
}
