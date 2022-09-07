// ==========================================================================
// 物联网快速开发自动生成model代码，无需手动修改，重新生成会自动覆盖.
// 生成日期：2022-06-21 22:06:17
// 生成路径: github.com/morgeq/iotfast/internal/app/model/entity/device_group.go
// 生成人：dwx
// ==========================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// DeviceGroup is the golang structure for table device_group.
type DeviceGroup struct {
	Id        int         `orm:"id,primary" json:"id"`        // 主键
	Name      string      `orm:"name" json:"name"`            // 组名
	CreatedAt *gtime.Time `orm:"created_at" json:"createdAt"` // 创建时间
	UpdatedAt *gtime.Time `orm:"updated_at" json:"updatedAt"` // 修改时间
}
