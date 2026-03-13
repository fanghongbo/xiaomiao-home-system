package utils

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// MySQL 的唯一键冲突错误代码是 1062，错误信息包含 "Duplicate entry"
func IsDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否是 MySQL 的唯一键冲突错误
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		// MySQL/MariaDB 唯一冲突错误码 1062
		return mysqlErr.Number == 1062
	}

	errStr := err.Error()
	// 检查错误信息中是否包含 "Duplicate entry"
	if strings.Contains(errStr, "Duplicate entry") {
		return true
	}

	// 检查是否是 GORM 的重复键错误
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// 检查 MySQL 错误代码 1062
	// MySQL 错误格式通常是: Error 1062: Duplicate entry 'xxx' for key 'xxx'
	if strings.Contains(errStr, "Error 1062") || strings.Contains(errStr, "1062") {
		return true
	}

	return false
}
