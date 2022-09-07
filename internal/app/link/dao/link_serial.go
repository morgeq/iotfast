// ==========================================================================
// 物联网快速开发自动生成dao操作代码，无需手动修改，重新生成不会自动覆盖.
// 生成日期：2022-07-16 15:59:37
// 生成路径: iotfast/internal/app/link/dao/link_serial.go
// 生成人：dwx
// ==========================================================================

package dao

import (
	"github.com/morgeq/iotfast/internal/app/link/dao/internal"
	// "github.com/gogf/gf/v2/os/gtime"
)

// internallinkSerialDao is internal type for wrapping internal DAO implements.
type internalLinkSerialDao = *internal.LinkSerialDao

// linkSerialDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type linkSerialDao struct {
	internalLinkSerialDao
}

var (
	// LinkSerial is globally public accessible object for table tools_gen_table operations.
	LinkSerial = linkSerialDao{
		internal.NewLinkSerialDao(),
	}
)
