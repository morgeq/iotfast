// ==========================================================================
// 物联网快速开发自动生成dao操作代码，无需手动修改，重新生成不会自动覆盖.
// 生成日期：2022-06-18 13:58:41
// 生成路径: github.com/morgeq/iotfast/internal/app/system/dao/sys_job.go
// 生成人：dwx
// ==========================================================================

package dao

import (
	"github.com/morgeq/iotfast/internal/app/system/dao/internal"
	// "github.com/gogf/gf/v2/os/gtime"
)

// internaljobDao is internal type for wrapping internal DAO implements.
type internalJobDao = *internal.SysJobDao
// jobDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type jobDao struct {
	internalJobDao
}
var (
    // SysJob is globally public accessible object for table tools_gen_table operations.
    SysJob = jobDao{
        internal.NewSysJobDao(),
    }
)
