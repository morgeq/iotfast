// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"github.com/morgeq/iotfast/internal/app/system/dao/internal"
)

// internalSysPostDao is internal type for wrapping internal DAO implements.
type internalSysPostDao = *internal.SysPostDao

// sysPostDao is the data access object for table sys_post.
// You can define custom methods on it to extend its functionality as you wish.
type sysPostDao struct {
	internalSysPostDao
}

var (
	// SysPost is globally public accessible object for table sys_post operations.
	SysPost = sysPostDao{
		internal.NewSysPostDao(),
	}
)

// Fill with you ideas below.
