// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"github.com/xiaodingding/iotfast/internal/app/system/dao/internal"
)

// internalSysLoginLogDao is internal type for wrapping internal DAO implements.
type internalSysLoginLogDao = *internal.SysLoginLogDao

// sysLoginLogDao is the data access object for table sys_login_log.
// You can define custom methods on it to extend its functionality as you wish.
type sysLoginLogDao struct {
	internalSysLoginLogDao
}

var (
	// SysLoginLog is globally public accessible object for table sys_login_log operations.
	SysLoginLog = sysLoginLogDao{
		internal.NewSysLoginLogDao(),
	}
)

// Fill with you ideas below.
