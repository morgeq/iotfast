// ==========================================================================
// 物联网快速开发自动生成dao操作代码，无需手动修改，重新生成不会自动覆盖.
// 生成日期：2022-06-21 22:06:17
// 生成路径: github.com/morgeq/iotfast/internal/app/device/dao/device_label_data.go
// 生成人：dwx
// ==========================================================================

package dao

import (
	"github.com/morgeq/iotfast/internal/app/device/dao/internal"
)

// internaldeviceLabelDataDao is internal type for wrapping internal DAO implements.
type internalDeviceLabelDataDao = *internal.DeviceLabelDataDao

// deviceLabelDataDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type deviceLabelDataDao struct {
	internalDeviceLabelDataDao
}

var (
	// DeviceLabelData is globally public accessible object for table tools_gen_table operations.
	DeviceLabelData = deviceLabelDataDao{
		internal.NewDeviceLabelDataDao(),
	}
)
